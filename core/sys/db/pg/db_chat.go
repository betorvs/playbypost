package pg

import (
	"context"
	"slices"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) AddChatInformation(ctx context.Context, username, userid, channel, chat string) (int, error) {
	query := "INSERT INTO chat_information(userid, channel, username, chat) VALUES($1, $2, $3, $4) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare insert into chat information failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(userid, channel, username, chat).Scan(&res)
	if err != nil {
		db.Logger.Error("query row insert into chat information failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) GetChatInformation(ctx context.Context) ([]types.ChatInfo, error) {
	info := []types.ChatInfo{}
	infoMap := make(map[string]types.ChatInfo)
	query := "SELECT id, userid, username, channel, chat FROM chat_information"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on chat_information failed", "error", err.Error())
		return info, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.ChatInfo
		if err := rows.Scan(&s.ID, &s.UserID, &s.Username, &s.Channel, &s.Chat); err != nil {
			db.Logger.Error("scan error on chat_information", "error", err.Error())
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
		db.Logger.Error("rows err on chat_information", "error", err.Error())
	}
	for _, v := range infoMap {
		info = append(info, v)
	}
	return info, nil
}

func (db *DBX) GetChatChannelInformation(ctx context.Context) ([]string, error) {
	info := []string{}
	query := "SELECT channel FROM chat_information"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on chat_information.channel failed", "error", err.Error())
		return info, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			db.Logger.Error("scan error on chat_information.channel", "error", err.Error())
		}
		if !slices.Contains(info, s) {
			info = append(info, s)
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on chat_information.channel", "error", err.Error())
	}
	return info, nil
}

// func (db *DBX) GetUserCard(ctx context.Context) ([]types.ChatInfo, error) {

// }
