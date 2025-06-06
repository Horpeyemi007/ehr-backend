package main

import (
	"backend/ehr/internal/auth"
	"backend/ehr/internal/config"
	"backend/ehr/internal/db"
	"backend/ehr/internal/env"
	"backend/ehr/internal/logging"
	"backend/ehr/internal/model"
	"log"
	"time"

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
		Auth: config.AuthConfig{
			Token: config.TokenConfig{
				PrivateKeyPath: env.GetPrivateKeyPath(),
				PublicKeyPath:  env.GetPublicKeyPath(),
				Exp:            time.Hour * 24 * 3, // 3 days
				Iss:            "ehr",
			},
		},
	}
	// Initialize Global Logger
	logging.InitializeLogger()
	defer logging.CleanupLogger()

	pool, err := db.New(cfg.DB.Addr, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	// close the db connection
	defer pool.Close()
	log.Println("Database connection pool established")

	// initialize the db connection
	store := model.InitializeStore(pool)
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.Auth.Token.PrivateKeyPath,
		cfg.Auth.Token.PublicKeyPath,
		cfg.Auth.Token.Iss,
		cfg.Auth.Token.Iss,
	)
	app := newApplication(cfg, store, jwtAuthenticator)

	router := setupRouter(app)
	app.Run(router)

}
