package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/Sabirk8992/ecom-backend/internal/handler"
)

func Run(cfg *config.Config) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheck)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s [env=%s]", addr, cfg.AppEnv)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
