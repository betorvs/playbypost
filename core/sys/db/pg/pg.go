package pg

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type DBX struct {
	Conn   *sql.DB
	Logger *slog.Logger
}

func New(conn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
