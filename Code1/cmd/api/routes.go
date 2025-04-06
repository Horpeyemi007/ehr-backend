package main

import (
	"backend/ehr/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	Config  config.ServerConfig
	Patient patientHandler
}

func (app *application) Mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx)
	r.Use(middleware.Timeout(60 * time.Second))

	// Custom middleware for colored logging
	// r.Use(ColoredLogger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Route("/patient", func(r chi.Router) {
			r.Post("/create", app.Patient.CreatePatientHandler)
		})
	})
	return r
}

func ColoredLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		// Log with colors
		color.Set(color.FgGreen)
		defer color.Unset()
		color.New(color.FgCyan).Printf("Method: %s ", r.Method)
		color.New(color.FgYellow).Printf("URL: %s ", r.URL.Path)
		color.New(color.FgMagenta).Printf("Duration: %s\n", duration)
	})
}

// Run starts the server and it listen on the specified address
func (app *application) Run(mux http.Handler) error {
	server := &http.Server{
		Addr:         app.Config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.Config.Addr)

	return server.ListenAndServe()
}
