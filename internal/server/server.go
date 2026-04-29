package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/Sabirk8992/ecom-backend/internal/handler"
	"github.com/Sabirk8992/ecom-backend/internal/service"
)

func Run(cfg *config.Config, db *sql.DB) {
	authSvc := service.NewAuthService(db, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/auth/signup", authHandler.Signup)
	mux.HandleFunc("/auth/login", authHandler.Login)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s [env=%s]", addr, cfg.AppEnv)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
