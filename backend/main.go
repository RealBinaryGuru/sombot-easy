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
	"github.com/joho/godotenv"
	slogchi "github.com/samber/slog-chi"
)

// Constants
const uploadDir = "./uploads"

var promotionDir = filepath.Join(uploadDir, "promotions")

// Models
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
func respondJSON(w http.ResponseWriter, status int, response map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func successResponse(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"error":   false,
		"message": "success",
		"data":    data,
	})
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"error":   true,
		"message": message,
	})
}

// Utility Functions
func ensureDirectoryExists(dir string) error {
	dir = filepath.Clean(dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func generateMD5Hash(input string) string {
	hash := md5.New()
	_, _ = hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Database Connection
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

// Promotion Handlers
func createPromotionHandler(db *pg.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			errorResponse(w, http.StatusBadRequest, "Invalid form data: "+err.Error())
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "Failed to retrieve file: "+err.Error())
			return
		}
		defer file.Close()

		if err := ensureDirectoryExists(promotionDir); err != nil {
			errorResponse(w, http.StatusInternalServerError, "Failed to create directory: "+err.Error())
			return
		}

		ext := filepath.Ext(fileHeader.Filename)
		hashFilename := generateMD5Hash(time.Now().String() + fileHeader.Filename)
		filePath := filepath.Join(promotionDir, hashFilename+ext)

		outFile, err := os.Create(filepath.Clean(filePath))
		if err != nil {
			errorResponse(w, http.StatusInternalServerError, "Failed to save file: "+err.Error())
			return
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, file); err != nil {
			errorResponse(w, http.StatusInternalServerError, "Failed to write file: "+err.Error())
			return
		}

		promotion := &Promotion{
			PromotionName: fileHeader.Filename,
			ImageURL:      "/uploads/promotions/" + hashFilename + ext,
			StartDate:     time.Now(),
			EndDate:       time.Now().AddDate(0, 1, 0),
			Status:        "active",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		if _, err := db.Model(promotion).Insert(); err != nil {
			errorResponse(w, http.StatusInternalServerError, "Failed to create promotion: "+err.Error())
			return
		}

		successResponse(w, promotion)
	}
}

func getPromotionsHandler(db *pg.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var promotions []Promotion
		if err := db.Model(&promotions).Select(); err != nil {
			logger.Error("Error fetching promotions", slog.String("error", err.Error()))
			errorResponse(w, http.StatusInternalServerError, "Error fetching promotions")
			return
		}
		successResponse(w, promotions)
	}
}

// Setup Routes
func setupRoutes(r *chi.Mux, db *pg.DB, logger *slog.Logger) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.Route("/api/promotions", func(r chi.Router) {
		r.Post("/", createPromotionHandler(db, logger))
		r.Get("/", getPromotionsHandler(db, logger))
	})
	r.Get("/uploads/*", http.StripPrefix("/uploads", http.FileServer(http.Dir("./uploads"))).ServeHTTP)
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := connectToDatabase()
	if err != nil {
		logger.Error("Error connecting to database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(slogchi.New(logger))
	r.Use(middleware.Recoverer)

	if err := ensureDirectoryExists(promotionDir); err != nil {
		logger.Error("Failed to initialize uploads directory", slog.String("error", err.Error()))
		return
	}

	setupRoutes(r, db, logger)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	logger.Info("Starting server", slog.String("port", port))
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start", slog.String("error", err.Error()))
	}
}
