package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Sabirk8992/ecom-backend/internal/config"
	"github.com/Sabirk8992/ecom-backend/internal/handler"
	"github.com/Sabirk8992/ecom-backend/internal/logger"
	"github.com/Sabirk8992/ecom-backend/internal/metrics"
	"github.com/Sabirk8992/ecom-backend/internal/middleware"
	"github.com/Sabirk8992/ecom-backend/internal/service"
	"github.com/Sabirk8992/ecom-backend/internal/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
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
		logger.Log.Fatal("Failed to initialize S3", zap.Error(err))
	}
	uploadHandler := handler.NewUploadHandler(s3Storage)

	_ = metrics.HttpRequestsTotal

	mux := http.NewServeMux()

	// metrics endpoint for Prometheus to scrape
	mux.Handle("/metrics", promhttp.Handler())

	// health
	mux.HandleFunc("/health", middleware.Observability(handler.HealthCheck))

	// auth
	mux.HandleFunc("/auth/signup", middleware.Observability(authHandler.Signup))
	mux.HandleFunc("/auth/login", middleware.Observability(authHandler.Login))

	// products
	mux.HandleFunc("/products", middleware.Observability(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetAll(w, r)
		case http.MethodPost:
			productHandler.Create(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/products/", middleware.Observability(func(w http.ResponseWriter, r *http.Request) {
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
	}))

	// orders
	mux.HandleFunc("/orders", middleware.Observability(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(cfg.JWTSecret, orderHandler.Create)(w, r)
		case http.MethodGet:
			middleware.AuthMiddleware(cfg.JWTSecret, orderHandler.GetAll)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/orders/", middleware.Observability(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.AuthMiddleware(cfg.JWTSecret, orderHandler.GetByID)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// payments
	mux.HandleFunc("/payments", middleware.Observability(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(cfg.JWTSecret, paymentHandler.Process)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// upload
	mux.HandleFunc("/upload", middleware.Observability(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.AuthMiddleware(cfg.JWTSecret, uploadHandler.Upload)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	logger.Log.Info("Server starting", zap.String("addr", addr), zap.String("env", cfg.AppEnv))

	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}
