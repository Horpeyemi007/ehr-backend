package main

import (
	"backend/ehr/internal/config"
	"backend/ehr/internal/db"
	"backend/ehr/internal/env"
	"backend/ehr/internal/logging"
	"backend/ehr/internal/model"
	"backend/ehr/internal/service"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize env variables or fallback to default
	cfg := config.ServerConfig{
		Addr: env.GetString("ADDR", ":9000"),
		DB: config.DbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://postgres:adminuser@localhost/ehrdb?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.GetString("ENV", "dev"),
	}

	// Initialize Global Logger
	logging.InitializeLogger()
	defer logging.CleanupLogger()

	db, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	// close the db connection
	defer db.Close()
	log.Println("Database connection pool established")

	// initialize the model, services & handlers
	patientModel := model.NewPatientRepository(db)
	patientService := service.NewPatientService(patientModel)
	patientHandler := NewPatientHandler(patientService)

	app := &application{
		Config:  cfg,
		Patient: *patientHandler,
	}

	mux := app.Mount()
	log.Fatal(app.Run(mux))

}
