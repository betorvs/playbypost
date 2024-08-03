package pg

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(conn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
