package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/Sabirk8992/ecom-backend/internal/handler"
	"github.com/Sabirk8992/ecom-backend/internal/middleware"
	"github.com/Sabirk8992/ecom-backend/internal/service"
	"github.com/Sabirk8992/ecom-backend/internal/storage"
)

func Run(cfg *config.Config, db *sql.DB) {
	authSvc := service.NewAuthService(db, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authSvc)

	productSvc := service.NewProductService(db)
	productHandler := handler.NewProductHandler(productSvc)

	orderSvc := service.NewOrderService(db)
	orderHandler := handler.NewOrderHandler(orderSvc)

	paymentSvc := service.NewPaymentService(db)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	s3Storage, err := storage.NewS3Storage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3: %v", err)
	}
	uploadHandler := handler.NewUploadHandler(s3Storage)

	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/health", handler.HealthCheck)

	// auth
	mux.HandleFunc("/auth/signup", authHandler.Signup)
	mux.HandleFunc("/auth/login", authHandler.Login)

	// products
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetAll(w, r)
		case http.MethodPost:
			productHandler.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetByID(w, r)
		case http.MethodPut:
			productHandler.Update(w, r)
		case http.MethodDelete:
			productHandler.Delete(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// orders (protected)
	mux.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(cfg.JWTSecret, orderHandler.Create)(w, r)
		case http.MethodGet:
			middleware.AuthMiddleware(cfg.JWTSecret, orderHandler.GetAll)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.AuthMiddleware(cfg.JWTSecret, orderHandler.GetByID)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// payments (protected)
	mux.HandleFunc("/payments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(cfg.JWTSecret, paymentHandler.Process)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// upload (protected)
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(cfg.JWTSecret, uploadHandler.Upload)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s [env=%s]", addr, cfg.AppEnv)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
