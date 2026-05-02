package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Sabirk8992/ecom-backend/internal/logger"
	"github.com/Sabirk8992/ecom-backend/internal/metrics"
	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Observability(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next(wrapped, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(wrapped.status)

		// Record metrics
		metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)

		// Structured log
		logger.Log.Info("request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("status", status),
			zap.Float64("duration_seconds", duration),
		)
	}
}
