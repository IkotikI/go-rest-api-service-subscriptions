package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/config"
	"github.com/ikotiki/go-rest-api-service-subscriptions/internal/storage/postgresql"
	"github.com/ikotiki/go-rest-api-service-subscriptions/logger"

	_ "github.com/joho/godotenv/autoload"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.MustLoad()
	log := logger.InitLoggerByFlag("trace", true)

	db := MustMakeStorage(cfg)

	var err error
	switch os.Args[1] {
	case "up":
		err = goose.Up(db, cfg.DB.MigrationDir)
	case "down":
		err = goose.Down(db, cfg.DB.MigrationDir)
	default:
		log.Fatal().Msgf("Available commands: up, down: %s, provided: %s", os.Args[1])
	}
	if err != nil {
		log.Fatal().Msgf("Error running migration: %s", err)
	}

}

func MustMakeStorage(cfg *config.Config) *sql.DB {

	db, err := postgresql.NewSQLStorage(postgresql.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBname,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}

	return db.SQLInstance()

}
