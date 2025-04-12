package main

import (
	"backend/ehr/internal/config"
	"backend/ehr/internal/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type application struct {
	config config.ServerConfig
	store  model.Storage
}

func newApplication(cfg config.ServerConfig, store model.Storage) *application {
	return &application{
		store:  store,
		config: cfg,
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

	log.Printf("Server running on %s", app.config.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
