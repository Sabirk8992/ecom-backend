package service

import (
	"database/sql"
	"errors"

	"github.com/Sabirk8992/ecom-backend/internal/model"
)

type PaymentService struct {
	DB *sql.DB
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{DB: db}
}

func (s *PaymentService) Process(userID int, req model.PaymentRequest) (*model.PaymentResponse, error) {
	// 1. Fetch the order
	var order model.Order
	err := s.DB.QueryRow(
		`SELECT id, user_id, total, status FROM orders WHERE id = $1`,
		req.OrderID,
	).Scan(&order.ID, &order.UserID, &order.Total, &order.Status)

	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	if err != nil {
		return nil, err
	}

	// 2. Only the order owner can pay
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// 3. Prevent double payment
	if order.Status == "paid" {
		return nil, errors.New("order already paid")
	}
	if order.Status == "failed" {
		return nil, errors.New("order payment already failed")
	}

	// 4. Simulate payment result
	paymentStatus := "success"
	orderStatus := "paid"
	message := "payment successful"

	if req.SimulateFailure {
		paymentStatus = "failed"
		orderStatus = "failed"
		message = "payment failed - simulated failure"
	}

	// 5. Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 6. Insert payment record
	var payment model.Payment
	err = tx.QueryRow(
		`INSERT INTO payments (order_id, amount, status, method)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, order_id, amount, status, method, created_at`,
		req.OrderID, order.Total, paymentStatus, req.Method,
	).Scan(&payment.ID, &payment.OrderID, &payment.Amount, &payment.Status, &payment.Method, &payment.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 7. Update order status
	_, err = tx.Exec(
		`UPDATE orders SET status = $1 WHERE id = $2`,
		orderStatus, req.OrderID,
	)
	if err != nil {
		return nil, err
	}

	// 8. Commit
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.PaymentResponse{
		Payment: &payment,
		Message: message,
	}, nil
}
