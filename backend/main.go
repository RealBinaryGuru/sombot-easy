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
	"github.com/joho/godotenv"

	slogchi "github.com/samber/slog-chi"
)

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date string `json:"date"`
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "success",
		"data":    data,
	})
}

func SuccessPaginationResponse(w http.ResponseWriter, data interface{}, page, size, total int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "success",
		"data":    data,
		"page":    page,
		"size":    size,
		"total":   total,
	})
}

func BadRequest(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "data validation failed",
		"data":    data,
	})
}

func InternalServerError(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "internal server error",
		"data":    data,
	})
}

func connectToDatabase() (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT")),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	})

	_, err := db.Exec("SELECT 1")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env", slog.String("error", err.Error()))
	}

	// Connect to database
	db, err := connectToDatabase()
	if err != nil {
		logger.Error("Error connecting to database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(slogchi.New(logger))
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, World!"))
		})

		// Events endpoint
		r.Route("/events", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				var events []Event
				err := db.Model(&events).Select()
				if err != nil {
					logger.Error("Error fetching events", slog.String("error", err.Error()))
					InternalServerError(w, err.Error())
					return
				}

				SuccessResponse(w, events)
			})
		})
	})

	// Start the server
	port := os.Getenv("SERVER_PORT")
	logger.Info("Starting server", slog.String("port", port))

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start", slog.String("error", err.Error()))
	}
}
