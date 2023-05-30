package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi"
	"golang.org/x/exp/slog"
)

func main() {
	port := "8080"
	appEnv := "dev"

	if fromEnv := os.Getenv("ENV"); fromEnv != "" {
		appEnv = fromEnv
	}

	// Resources for logging
	// https://betterstack.com/community/guides/logging/logging-in-go/

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug, // we should toggle this if we're in prod
	}

	var handler slog.Handler = slog.NewTextHandler(os.Stdout, opts)
	if appEnv == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger) // Set the default logger >:)

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	logger.Info("Starting server...", "SERVER", fmt.Sprintf("http://localhost:%s", port))

	r := chi.NewRouter()

	r.HandleFunc("/{username}", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Called generic handler", "method", r.Method, "params", r.URL.Query())
		username := chi.URLParam(r, "username") // 👈 get param
		json.NewEncoder(w).Encode(struct {
			Key   string
			Value string
			Query url.Values
		}{
			"username",
			username,
			r.URL.Query(),
		})
	})

	// Setup filehandling
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	// File handling

	r.Options("/htmx/", func(w http.ResponseWriter, r *http.Request) {
		slog.Error("Called without implementation", "method", r.Method)
	})

	r.Get("/htmx/", func(w http.ResponseWriter, r *http.Request) {
		slog.Error("Called without implementation", "method", r.Method)

	})
	r.Put("/htmx/", func(w http.ResponseWriter, r *http.Request) {
		slog.Error("Called without implementation", "method", r.Method)

	})
	r.Post("/htmx/", func(w http.ResponseWriter, r *http.Request) {
		slog.Error("Called without implementation", "method", r.Method)
		fmt.Fprintf(w, "<strong>This is plain HTML</strong>")
	})

	r.Patch("/htmx/", func(w http.ResponseWriter, r *http.Request) {
		slog.Error("Called without implementation", "method", r.Method)

	})
	r.Delete("/htmx/", func(w http.ResponseWriter, r *http.Request) {
		slog.Error("Called without implementation", "method", r.Method)

	})

	// Apply auth middleware to only `GET /users/{id}`
	// router.Group(func(r chi.Router) {
	// 	r.Use(AuthMiddleware)
	// 	r.Get("/users/{id}")
	// })
	// r.Mount("/posts", postsResource{}.Routes())

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
