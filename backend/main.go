package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"

	slogchi "github.com/samber/slog-chi"
)

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

type RequestEventBody struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
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

	r.Route("/api/events", func(r chi.Router) {
		r.Get("/", GetEventsHandler(db, logger))
		r.Post("/", CreateEventHandler(db, logger, validate))
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
