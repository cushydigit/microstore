package repository

import (
	"database/sql"
	"errors"

	"github.com/cushydigit/microstore/shared/types"
)

type PostgresProductRepo struct {
	DB *sql.DB
}

func NewPostgresProductRepo(db *sql.DB) *PostgresProductRepo {
	return &PostgresProductRepo{DB: db}
}

func (r *PostgresProductRepo) Create(p *types.Product) error {
	return r.DB.QueryRow(
		`INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id`,
		p.Name,
		p.Description,
		p.Price,
		p.Stock,
	).Scan(&p.ID)
}

func (r *PostgresProductRepo) GetAll() ([]types.Product, error) {
	rows, err := r.DB.Query(`SELECT id, name, description, price, stock FROM products`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		var p types.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *PostgresProductRepo) GetByID(id int64) (*types.Product, error) {
	var p types.Product
	err := r.DB.QueryRow(`SELECT id, name, description, price, stock FROM products WHERE id = $1`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PostgresProductRepo) Delete(id int64) error {
	result, err := r.DB.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}
