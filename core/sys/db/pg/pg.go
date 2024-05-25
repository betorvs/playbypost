package pg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func New(conn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
