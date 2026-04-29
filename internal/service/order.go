package service

import (
	"database/sql"
	"errors"

	"github.com/Sabirk8992/ecom-backend/internal/model"
)

type OrderService struct {
	DB *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
	return &OrderService{DB: db}
}

func (s *OrderService) Create(userID int, req model.CreateOrderRequest) (*model.Order, error) {
	// Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // rollback if anything fails

	// 1. Lock the product row + check stock
	var price float64
	var stock int
	err = tx.QueryRow(
		`SELECT price, stock FROM products WHERE id = $1 FOR UPDATE`,
		req.ProductID,
	).Scan(&price, &stock)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	// 2. Check enough stock
	if stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// 3. Deduct stock
	_, err = tx.Exec(
		`UPDATE products SET stock = stock - $1 WHERE id = $2`,
		req.Quantity, req.ProductID,
	)
	if err != nil {
		return nil, err
	}

	// 4. Create order
	total := price * float64(req.Quantity)
	var order model.Order
	err = tx.QueryRow(
		`INSERT INTO orders (user_id, product_id, quantity, total, status)
		 VALUES ($1, $2, $3, $4, 'confirmed')
		 RETURNING id, user_id, product_id, quantity, total, status, created_at`,
		userID, req.ProductID, req.Quantity, total,
	).Scan(&order.ID, &order.UserID, &order.ProductID, &order.Quantity, &order.Total, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 5. Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrderService) GetAll() ([]model.Order, error) {
	rows, err := s.DB.Query(
		`SELECT id, user_id, product_id, quantity, total, status, created_at
		 FROM orders ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Total, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (s *OrderService) GetByID(id int) (*model.Order, error) {
	var o model.Order
	err := s.DB.QueryRow(
		`SELECT id, user_id, product_id, quantity, total, status, created_at
		 FROM orders WHERE id = $1`, id,
	).Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Total, &o.Status, &o.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}
	return &o, err
}
