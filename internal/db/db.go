package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Sabirk8992/ecom-backend/internal/config"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to reach DB: %v", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db
}