package repository

import (
	"database/sql"
	"errors"
	"github.com/cushydigit/microstore/shared/types"
)

type PostgresUserRepo struct {
	DB *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{DB: db}
}

func (r *PostgresUserRepo) FindByEmail(email string) (*types.User, error) {
	row := r.DB.QueryRow(
		`SELECT id, email, password FROM users WHERE email = $1`,
		email,
	)

	var user types.User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepo) Create(user *types.User) error {
	// check if the user already exists by email
	existingUser, _ := r.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	// Insert new user into the database
	err := r.DB.QueryRow(
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		user.Email,
		user.Password,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}
