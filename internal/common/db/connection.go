package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(connString string) (*sql.DB, error) {

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
