package main

import (
	"backend/ehr/cmd/api"
	"backend/ehr/internal/config"
	"backend/ehr/internal/db"
	"backend/ehr/internal/env"
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
		Addr: env.GetString("ADDR", ":9090"),
		DB: config.DbConfig{
			Addr:         env.GetString("DB_ADDR", "admin:adminuser@tcp(localhost:3306)/eazybank"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Env: env.GetString("ENV", "dev"),
	}

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

	app := &api.Application{
		Config: cfg,
	}

	mux := app.Mount()
	log.Fatal(app.Run(mux))

}
