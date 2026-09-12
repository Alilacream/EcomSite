package store

import (
	"context"
	"database/sql"
	"errors"

	"alilacream/ecom/internal/models"
	"alilacream/ecom/lib"
)

type CustomerRepository interface {
	GetbyID(ctx context.Context, id string) (*models.Customer, error)
	VerifyCrendiatials(ctx context.Context, customer *models.Customer) (*models.Customer, error)
	GetbyUsername(ctx context.Context, username string) (customer *models.Customer, err error)
	Create(ctx context.Context, customer *models.Customer) (*models.Customer, error)
}

type CustomerStore struct {
	db *sql.DB
}

// @GET by id
// Utile for Admin dashboard
func (s *CustomerStore) GetbyID(ctx context.Context, id string) (customer *models.Customer, err error) {
	query := `SELECT username, email, created_at WHERE id = $1`
	err = s.db.QueryRowContext(ctx, query, id).Scan(&customer.Username, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// @GET By username
func (s *CustomerStore) GetbyUsername(ctx context.Context, username string) (customer *models.Customer, err error) {
	query := `SELECT username, email, created_at WHERE username = $1`
	err = s.db.QueryRowContext(ctx, query, username).Scan(&customer.Username, &customer.Email, &customer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

// @CREATE
func (s *CustomerStore) Create(ctx context.Context, newCustomer *models.Customer) (*models.Customer, error) {
	hashedPass, err := lib.HashPassword(newCustomer.Password)
	if err != nil {
		return nil, err
	}
	query := `INSERT INTO customers (username, email, password) VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`
	error := s.db.QueryRowContext(ctx, query, newCustomer.Username, newCustomer.Email, hashedPass).Scan(
		&newCustomer.ID,
		&newCustomer.Username,
		&newCustomer.Email,
		&newCustomer.CreatedAt,
	)
	return newCustomer, error
}

// @GET & VERIFY
// this method is utile for our login func
func (s *CustomerStore) VerifyCrendiatials(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	var storedClient models.Customer
	query := `SELECT id, email, username, created_at , password FROM customers WHERE username = $1`
	err := s.db.QueryRowContext(ctx, query, customer.Username).Scan(
		&storedClient.ID,
		&storedClient.Email,
		&storedClient.Username,
		&storedClient.CreatedAt,
		&storedClient.Password,
	)
	if err != nil {
		return nil, err
	}
	if ok := lib.CheckPasswordHash(customer.Password, storedClient.Password); !ok {
		return nil, errors.New("Invalid Password")
	}
	return &storedClient, nil
}
