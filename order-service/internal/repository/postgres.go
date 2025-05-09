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
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Insert order
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, total_price, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`, order.UserID, order.TotalPrice, order.Status).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert order items
	for _, item := range order.Items {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity)
			VALUES ($1, $2, $3)
		`, order.ID, item.ProductID, item.Quantity)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) GetByID(id int64) (*types.Order, error) {
	order := &types.Order{ID: id}
	err := r.DB.QueryRow(`
		SELECT user_id, total_price, status FROM orders WHERE id = $1
	`, id).Scan(&order.UserID, &order.TotalPrice, &order.Status)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(`
		SELECT product_id, quantity FROM order_items WHERE order_id = $1
	`, id)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var item types.OrderItem
		if err := rows.Scan(&item.ProductID, &item.Quantity); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (r *PostgresOrderRepository) GetByUserID(userID int) ([]types.Order, error) {
	rows, err := r.DB.Query(`
		SELECT id FROM orders WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}

	var orders []types.Order
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		order, err := r.GetByID(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}

func (r *PostgresOrderRepository) GetAll() ([]types.Order, error) {
	rows, err := r.DB.Query(`
		SELECT id FROM orders
	`)
	if err != nil {
		return nil, err
	}
	var orders []types.Order
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		order, err := r.GetByID(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}
