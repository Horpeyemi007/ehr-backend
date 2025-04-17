package main

import (
	"backend/ehr/internal/auth"
	"backend/ehr/internal/config"
	"backend/ehr/internal/logging"
	"backend/ehr/internal/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type application struct {
	config        config.ServerConfig
	store         model.Storage
	authenticator auth.Authenticator
}

func newApplication(
	cfg config.ServerConfig,
	store model.Storage,
	auth auth.Authenticator,
) *application {
	return &application{
		store:         store,
		config:        cfg,
		authenticator: auth,
	}
}

func (app *application) Run(router *gin.Engine) {
	s := &http.Server{
		Addr:         app.config.Addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	logging.Logger.Infow("Server running on port", "addr", app.config.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
