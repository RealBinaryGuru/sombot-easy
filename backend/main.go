package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Set up the router
	r := chi.NewRouter()
	r.Use(slogchi.New(logger)) // Middleware to use slog for request logging
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(400), 400)
	})

	port := ":3000"
	logger.Info("Starting server", slog.String("port", port))

	if err := http.ListenAndServe(port, r); err != nil {
		logger.Error("Server failed to start", slog.String("error", err.Error()))
	}
}
