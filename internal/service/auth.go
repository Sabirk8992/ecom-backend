package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Sabirk8992/ecom-backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB        *sql.DB
	JWTSecret string
}

func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{DB: db, JWTSecret: jwtSecret}
}

func (s *AuthService) Signup(req model.SignupRequest) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(
		"INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		req.Name, req.Email, string(hashed),
	)
	return err
}

func (s *AuthService) Login(req model.LoginRequest) (string, error) {
	var user model.User

	err := s.DB.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return "", errors.New("invalid credentials")
	}
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(s.JWTSecret))
}
