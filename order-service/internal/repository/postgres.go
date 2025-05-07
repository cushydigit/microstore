package repository

import (
	"database/sql"

	"github.com/cushydigit/microstore/shared/types"
)

type PostgresOrderRepository struct {
	DB *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db}
}

func (r *PostgresOrderRepository) Create(order *types.Order) error {

	return nil
}

func (r *PostgresOrderRepository) GetByID(id int64) (*types.Order, error) {
	return nil, nil
}

func (r *PostgresOrderRepository) GetByUserID(userID int) ([]*types.Order, error) {
	return nil, nil
}

func (r *PostgresOrderRepository) GetAll() ([]*types.Order, error) {
	return nil, nil
}
