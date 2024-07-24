package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shammianand/fast-json-viewer/internal/handlers"
	"github.com/shammianand/fast-json-viewer/internal/middleware"
	"github.com/shammianand/fast-json-viewer/internal/services"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("web/templates/*.html"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	templates.ExecuteTemplate(w, "layout.html", nil)
}

func main() {
	logger := log.New(os.Stdout, "[JSON-VIEWER-LOGS]: ", log.LstdFlags)

	sessionManager := services.NewSessionManager()
	parser := services.NewParser(300 * 1024 * 1024) // 300 MB max file size

	mux := http.NewServeMux()

	mux.Handle("/", middleware.Logging(logger)(
		middleware.RateLimit(100, time.Minute)(
			http.HandlerFunc(homeHandler),
		),
	))

	mux.Handle("/upload/", middleware.Logging(logger)(
		middleware.RateLimit(10, time.Minute)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlers.UploadHandler(w, r, sessionManager, parser)
			}),
		),
	))

	mux.Handle("/structure/", middleware.Logging(logger)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.GetStructureHandler(w, r, sessionManager)
		}),
	))

	mux.Handle("/expand/", middleware.Logging(logger)(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.ExpandNodeHandler(w, r, sessionManager)
		}),
	))

	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
