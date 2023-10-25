package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slog"
)

var (
	logger *slog.Logger
)

func main() {
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

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
	logger = slog.New(handler)

	logger.Info("Starting server...", "server", fmt.Sprintf("http://localhost:%s", port))

	r := mux.NewRouter()

	// Set no caching
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
			wr.Header().Set("Cache-Control", "max-age=0, must-revalidate")
			next.ServeHTTP(wr, req)
		})
	})

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseGlob("partials/*.tpl.html"))
		// You can see the available templates here
		// log.Println(tmpl.DefinedTemplates())
		err := tmpl.Lookup("index.tpl.html").Execute(w, nil)
		if err != nil {
			logger.Error("Failed to execute template", "template", tmpl.Name, "error", err)
		}
	})

	r.HandleFunc("/grab-it", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseGlob("partials/*.tpl.html"))
		err := tmpl.Lookup("includeme2").Execute(w, nil)
		if err != nil {
			logger.Error("Failed to execute template", "template", tmpl.Name, "error", err)
		}
	})

	// Setup filehandling
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))

	r.HandleFunc("/breakfast", func(w http.ResponseWriter, r *http.Request) {
		//slog.Error("Called without implementation", "method", r.Method)
		partialEncoder(w, "results", map[string]string{
			"Eggs":     "Fried",
			"Bacon":    "Venison",
			"Sausages": "Pork",
		})
	})

	r.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Called search handler", "method", r.Method, "params", r.URL.Query())
		query := r.Form.Get("query")
		json.NewEncoder(w).Encode(struct {
			Key   string
			Value string
			Query url.Values
		}{
			"query",
			query,
			r.URL.Query(),
		})
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
