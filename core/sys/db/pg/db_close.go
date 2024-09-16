package pg

func (db *DBX) Close() error {
	return db.Conn.Close()
}
