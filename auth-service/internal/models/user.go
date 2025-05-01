package models

// TODO change ID int to int64
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
