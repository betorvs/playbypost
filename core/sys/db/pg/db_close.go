package pg

import (
	"fmt"

	"github.com/lib/pq"
)

func (db *DBX) Close() error {
	return db.Conn.Close()
}

func (db *DBX) parsePostgresError(err error) error {
	pgErr, ok := err.(*pq.Error)
	if ok {
		db.Logger.Error("pq error", "code", pgErr.Code, "message", pgErr.Message)
		return fmt.Errorf("db: %s", pgErr.Message)
	}
	return err
}
