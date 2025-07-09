package pg

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) CreateWriters(ctx context.Context, username, password string) (int, error) {
	query := "INSERT INTO writers(username, password) VALUES($1, $2) RETURNING id" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare insert into writers failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			db.Logger.Error("error closing stmt", "error", err)
		}
	}()
	var res int
	err = stmt.QueryRow(username, password).Scan(&res)
	if err != nil {
		db.Logger.Error("query row insert into writers failed", "error", err.Error())
		return -1, db.parsePostgresError(err)
	}
	return res, nil
}

func (db *DBX) GetWriters(ctx context.Context, active bool) ([]types.Writer, error) {
	query := "SELECT id, username FROM writers" // dev:finder+query
	if active {
		query = "SELECT id, username FROM writers WHERE active = true" // dev:finder+query
	}
	users := []types.Writer{}
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on writers failed", "error", err.Error())
		return users, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		var user types.Writer
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			db.Logger.Error("scan error on writers", "error", err.Error())
		}
		users = append(users, user)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on writers", "error", err.Error())
	}
	return users, nil
}

func (db *DBX) GetWriterByID(ctx context.Context, id int) (types.Writer, error) {
	user := types.Writer{}
	keys := make(map[int]string)
	query := "SELECT w.id, w.username, k.story_id, k.encoding_key FROM writers AS w JOIN access_story_keys AS a ON w.id = a.writer_id JOIN story_keys AS k ON a.story_keys_id = k.id WHERE w.id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on writers by id failed", "error", err.Error())
		return user, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		// var user types.Writer
		var key string
		var keyID int
		if err := rows.Scan(&user.ID, &user.Username, &keyID, &key); err != nil {
			db.Logger.Error("scan error on writers by id", "error", err.Error())
		}
		_, ok := keys[keyID]
		if !ok {
			keys[keyID] = key
		}
		// users = append(users, user)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on writers by id", "error", err.Error())
	}
	user.EncodingKeys = keys
	return user, nil
}

func (db *DBX) GetWriterByUsername(ctx context.Context, username string) (types.Writer, error) {
	user := types.Writer{}
	query := "SELECT id, username, password FROM writers WHERE username = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, username)
	if err != nil {
		db.Logger.Error("query on writers by username failed", "error", err.Error())
		return user, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		// var user types.Writer
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			db.Logger.Error("scan error on writers by username", "error", err.Error())
		}
		// users = append(users, user)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on writers by username", "error", err.Error())
	}
	return user, nil
}
