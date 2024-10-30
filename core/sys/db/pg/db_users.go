package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetUser(ctx context.Context) ([]types.User, error) {
	users := []types.User{}
	query := "SELECT id, userid, active FROM users" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on users failed", "error", err.Error())
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var u types.User
		if err := rows.Scan(&u.ID, &u.UserID, &u.Active); err != nil {
			db.Logger.Error("scan error on users", "error", err.Error())
		}
		users = append(users, u)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on users", "error", err.Error())
	}
	return users, nil
}

func (db *DBX) GetUserByUserID(ctx context.Context, id string) (types.User, error) {
	user := types.User{}
	query := "SELECT id, userid, active FROM users WHERE userid = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on users by userid failed", "error", err.Error())
		return user, err
	}
	defer rows.Close()
	for rows.Next() {
		var u types.User
		if err := rows.Scan(&u.ID, &u.UserID, &u.Active); err != nil {
			db.Logger.Error("scan error on users by userid", "error", err.Error())
		}
		if u.ID > 0 {
			user = u
		}
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on users by userid", "error", err.Error())
	}
	return user, nil
}

func (db *DBX) CreateUserTx(ctx context.Context, userid string) (int, error) {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on CreateUserTx failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// check user exist
	queryUser := "SELECT id FROM users WHERE userid = $1" // dev:finder+query
	stmtQueryUser, err := db.Conn.PrepareContext(ctx, queryUser)
	if err != nil {
		db.Logger.Error("tx prepare on queryUser failed", "error", err.Error())
		return -1, err
	}
	defer stmtQueryUser.Close()
	var userID int
	err = tx.StmtContext(ctx, stmtQueryUser).QueryRow(userid).Scan(&userID)
	if err != nil {
		db.Logger.Info("user not found creating new", "return", err.Error())
		// just log this error
		// return -1, err

	}
	if userID == 0 {
		id, err := db.createUser(ctx, userid, tx)
		if err != nil {
			db.Logger.Error("insert into users failed", "error", err.Error())
			return -1, err
		}
		userID = id
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on CreateUserTx failed", "error", err.Error())
		return -1, err
	}
	return userID, nil

}

func (db *DBX) createUser(ctx context.Context, userid string, tx *sql.Tx) (int, error) {
	var userID int
	queryInsertUser := "INSERT INTO users(userid) VALUES($1) RETURNING id" // dev:finder+query
	stmtInsertUser, err := db.Conn.PrepareContext(ctx, queryInsertUser)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer stmtInsertUser.Close()
	err = tx.StmtContext(ctx, stmtInsertUser).QueryRow(userid).Scan(&userID)
	if err != nil {
		db.Logger.Error("query row insert into users failed", "error", err.Error())
		return -1, err
	}
	return userID, nil
}
