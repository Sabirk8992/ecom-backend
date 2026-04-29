package service

import (
	"database/sql"
	"errors"

	"github.com/Sabirk8992/ecom-backend/internal/model"
)

type ProductService struct {
	DB *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) Create(req model.CreateProductRequest) (*model.Product, error) {
	var p model.Product
	err := s.DB.QueryRow(
		`INSERT INTO products (name, description, price, stock)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, description, price, stock, created_at`,
		req.Name, req.Description, req.Price, req.Stock,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)
	return &p, err
}

func (s *ProductService) GetAll() ([]model.Product, error) {
	rows, err := s.DB.Query(
		`SELECT id, name, description, price, stock, created_at FROM products ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *ProductService) GetByID(id int) (*model.Product, error) {
	var p model.Product
	err := s.DB.QueryRow(
		`SELECT id, name, description, price, stock, created_at FROM products WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return &p, err
}

func (s *ProductService) Update(id int, req model.CreateProductRequest) (*model.Product, error) {
	var p model.Product
	err := s.DB.QueryRow(
		`UPDATE products SET name=$1, description=$2, price=$3, stock=$4
		 WHERE id=$5
		 RETURNING id, name, description, price, stock, created_at`,
		req.Name, req.Description, req.Price, req.Stock, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	return &p, err
}

func (s *ProductService) Delete(id int) error {
	result, err := s.DB.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}
