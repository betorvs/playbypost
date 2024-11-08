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

func convertInterfaceInt(x interface{}) int {
	switch v := x.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		return 0
	default:
		return 0
	}
}
