package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	slogchi "github.com/samber/slog-chi"
)

// Events
type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

type RequestEventBody struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

// Promotions
type Promotion struct {
	PromotionID   int       `json:"promotion_id"`
	PromotionName string    `json:"promotion_name"`
	ImageURL      string    `json:"image_url"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Response Helpers
func respondJSON(w http.ResponseWriter, status int, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"error":   false,
		"message": "success",
		"data":    data,
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	respondJSON(w, http.StatusBadRequest, map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

func InternalServerError(w http.ResponseWriter, message string) {
	respondJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

// Function to get file extension
func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ".bin" // Default extension if none found
	}
	return ext
}

// Function to generate MD5 hash from string
func generateMD5Hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func connectToDatabase() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	})

	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, err
	}

	return db, nil
}

func GetEventsHandler(db *pg.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var events []Event
		if err := db.Model(&events).Select(); err != nil {
			logger.Error("Error fetching events", slog.String("error", err.Error()))
			InternalServerError(w, "Error fetching events")
			return
		}
		SuccessResponse(w, events)
	}
}

func CreateEventHandler(db *pg.DB, logger *slog.Logger, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData RequestEventBody

		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			BadRequest(w, "Invalid JSON payload")
			return
		}
		if err := validate.Struct(requestData); err != nil {
			BadRequest(w, "Validation error: "+err.Error())
			return
		}

		event := &Event{Name: requestData.Name}
		if _, err := db.Model(event).Insert(); err != nil {
			logger.Error("Error inserting event", slog.String("error", err.Error()))
			InternalServerError(w, "Failed to create event")
			return
		}
		SuccessResponse(w, event)
	}
}

func setupRoutes(r *chi.Mux, db *pg.DB, logger *slog.Logger, validate *validator.Validate) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/events", func(r chi.Router) {
			r.Get("/", GetEventsHandler(db, logger))
			r.Post("/", CreateEventHandler(db, logger, validate))
		})

		r.Route("/promotions", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				// Parse the multipart form (maximum upload size 10MB)
				err := r.ParseMultipartForm(10 << 20) // 10 MB
				if err != nil {
					BadRequest(w, err.Error())
					return
				}

				file, fileHeader, err := r.FormFile("file")
				if err != nil {
					BadRequest(w, err.Error())
					return
				}
				defer file.Close()

				// Ensure the uploads directory exists
				destDir := "./uploads"
				if _, err := os.Stat(destDir); os.IsNotExist(err) {
					err := os.MkdirAll(destDir, 0755) // Create the directory if it doesn't exist
					if err != nil {
						InternalServerError(w, err.Error())
						return
					}
				}

				// Get the file extension
				ext := getFileExtension(fileHeader.Filename)

				// Concatenate the current date and original filename to create a unique string
				dateStr := time.Now().Format("20060102_150405") // Example: "20231128_102030"
				uniqueString := dateStr + "_" + fileHeader.Filename

				// Generate the MD5 hash of the unique string
				hashFilename := generateMD5Hash(uniqueString)

				promotionsDir := filepath.Join(destDir, "promotions")

				// Ensure that the "promotions" directory exists
				if err := os.MkdirAll(promotionsDir, os.ModePerm); err != nil {
					InternalServerError(w, err.Error())
					return
				}

				// Join the "promotions" directory with the filename
				promotionFilePath := filepath.Join(promotionsDir, hashFilename+ext)

				outFile, err := os.Create(promotionFilePath)
				if err != nil {
					InternalServerError(w, err.Error())
					return
				}
				defer outFile.Close()

				// Copy the uploaded file to the destination
				_, err = io.Copy(outFile, file)
				if err != nil {
					InternalServerError(w, err.Error())
					return
				}

				promotion := &Promotion{
					PromotionName: fileHeader.Filename,
					ImageURL:      promotionFilePath,
					StartDate:     time.Now(),
					EndDate:       time.Now().AddDate(0, 1, 0),
					Status:        "active",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
				}

				if _, err := db.Model(promotion).Insert(); err != nil {
					InternalServerError(w, "Failed to create promotion")
					return
				}

				SuccessResponse(w, promotion)
			})
		})

	})

}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Connect to the database
	db, err := connectToDatabase()
	if err != nil {
		logger.Error("Error connecting to database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	// Initialize router and middlewares
	r := chi.NewRouter()
	r.Use(slogchi.New(logger))
	r.Use(middleware.Recoverer)
	// Serve static file
	r.Get("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads"))).ServeHTTP)

	// Setup routes
	validate := validator.New()
	setupRoutes(r, db, logger, validate)

	// Start the server
	port := os.Getenv("SERVER_PORT")
	logger.Info("Starting server", slog.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start", slog.String("error", err.Error()))
	}

}
