package db

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) AddSlackInformation(ctx context.Context, username, userid, channel string) (int, error) {
	query := "INSERT INTO slack_information(userid, channel, username) VALUES($1, $2, $3) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("prepare insert into slack information failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(userid, channel, username).Scan(&res)
	if err != nil {
		db.logger.Error("query row insert into slack information failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) GetSlackInformation(ctx context.Context) ([]types.SlackInfo, error) {
	info := []types.SlackInfo{}
	query := "SELECT id, userid, username, channel FROM slack_information"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.logger.Error("query on slack_information failed", "error", err.Error())
		return info, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.SlackInfo
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Channel); err != nil {
			db.logger.Error("scan error on slack_information", "error", err.Error())
		}
		info = append(info, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows err on slack_information", "error", err.Error())
	}
	return info, nil
}

// func (db *DBX) GetUserCard(ctx context.Context) ([]types.SlackInfo, error) {

// }
