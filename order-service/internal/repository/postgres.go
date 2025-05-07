package repository

import (
	"database/sql"

	"github.com/cushydigit/microstore/order-service/internal/models"
)

type PostgresOrderRepository struct {
	DB *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db}
}

func (r *PostgresOrderRepository) Create(order *models.Order) error {
	return nil
}

func (r *PostgresOrderRepository) GetByID(id int64) (*models.Order, error) {
	return nil, nil
}

func (r *PostgresOrderRepository) GetByUserID(userID int) ([]*models.Order, error) {
	return nil, nil
}

func (r *PostgresOrderRepository) GetAll() ([]*models.Order, error) {
	return nil, nil
}
