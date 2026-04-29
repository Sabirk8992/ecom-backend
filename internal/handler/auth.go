package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Sabirk8992/ecom-backend/internal/model"
	"github.com/Sabirk8992/ecom-backend/internal/service"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: svc}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var req model.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "name, email and password are required", http.StatusBadRequest)
		return
	}

	if err := h.Service.Signup(req); err != nil {
		log.Printf("Signup error: %v", err) // ← shows real error
		http.Error(w, "email already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user created successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(req)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.AuthResponse{Token: token})
}
