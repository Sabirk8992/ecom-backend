package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sabirk8992/ecom-backend/internal/middleware"
	"github.com/Sabirk8992/ecom-backend/internal/model"
	"github.com/Sabirk8992/ecom-backend/internal/service"
)

type PaymentHandler struct {
	Service *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{Service: svc}
}

func (h *PaymentHandler) Process(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.OrderID == 0 {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}

	if req.Method == "" {
		http.Error(w, "method is required (card/upi/wallet)", http.StatusBadRequest)
		return
	}

	resp, err := h.Service.Process(userID, req)
	if err != nil {
		log.Printf("Payment error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
