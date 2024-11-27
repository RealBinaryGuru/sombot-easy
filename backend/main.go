package main

import (
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error load .env", slog.String("error", err.Error()))
	}

	db, err := connectToDatabase()
	if err != nil {
		logger.Error("Error connecting to database", slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(slogchi.New(logger))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
	})

	port := os.Getenv("SERVER_PORT")
	logger.Info("Starting server", slog.String("port", port))

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logger.Error("Server failed to start", slog.String("error", err.Error()))
	}
}
