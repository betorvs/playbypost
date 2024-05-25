package db

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/betorvs/playbypost/core/sys/db/pg"
)

const (
	postgreSQL = "postgres"
	SQLite     = "sqlite"
)

type DBX struct {
	Conn   *sql.DB
	logger *slog.Logger
}

func NewDB(conn string, logger *slog.Logger) DBClient {
	db, err := pg.New(conn)
	if err != nil {
		logger.Error("sql open error", "error", err.Error())
		os.Exit(2)
	}

	// Set connection limits for connection pooling
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	// force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		logger.Error("error ping ", "error", err.Error())
		os.Exit(1)
	}
	logger.Info("connection to database okay")

	return &DBX{
		Conn:   db,
		logger: logger,
	}
}

func (db *DBX) Close() error {
	return db.Conn.Close()
}
