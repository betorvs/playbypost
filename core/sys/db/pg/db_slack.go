package pg

import (
	"context"
	"slices"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) AddSlackInformation(ctx context.Context, username, userid, channel string) (int, error) {
	query := "INSERT INTO slack_information(userid, channel, username) VALUES($1, $2, $3) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare insert into slack information failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(userid, channel, username).Scan(&res)
	if err != nil {
		db.Logger.Error("query row insert into slack information failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) GetSlackInformation(ctx context.Context) ([]types.SlackInfo, error) {
	info := []types.SlackInfo{}
	infoMap := make(map[string]types.SlackInfo)
	query := "SELECT id, userid, username, channel FROM slack_information"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on slack_information failed", "error", err.Error())
		return info, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.SlackInfo
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Channel); err != nil {
			db.Logger.Error("scan error on slack_information", "error", err.Error())
		}
		if v, ok := infoMap[s.UserID]; ok {
			db.Logger.Info("ok map", "values", v)
			if infoMap[s.UserID].Channel != s.Channel {
				// infoMap[s.UserID].Channel = infoMap[s.UserID].Channel + ", " + s.Channel
				s.Channel += ", " + infoMap[s.UserID].Channel
			}
			if infoMap[s.UserID].Username != s.Username {
				s.Username += ", " + infoMap[s.UserID].Username
			}
		}
		infoMap[s.UserID] = s
		// info = append(info, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on slack_information", "error", err.Error())
	}
	for _, v := range infoMap {
		info = append(info, v)
	}
	return info, nil
}

func (db *DBX) GetSlackChannelInformation(ctx context.Context) ([]string, error) {
	info := []string{}
	query := "SELECT channel FROM slack_information"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on slack_information.channel failed", "error", err.Error())
		return info, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			db.Logger.Error("scan error on slack_information.channel", "error", err.Error())
		}
		if !slices.Contains(info, s) {
			info = append(info, s)
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on slack_information.channel", "error", err.Error())
	}
	return info, nil
}

// func (db *DBX) GetUserCard(ctx context.Context) ([]types.SlackInfo, error) {

// }
