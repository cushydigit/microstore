
package database

import (
	"database/sql"
	"os"
	_ "github.com/lib/pq"

)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}
	return db, nil
}
