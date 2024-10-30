package pg

import (
	"context"
	"fmt"
	"slices"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) AddChatInformation(ctx context.Context, username, userid, channel, chat string) (int, error) {
	query := "INSERT INTO chat_information(userid, channel, username, chat) VALUES($1, $2, $3, $4) RETURNING id" // dev:finder+query
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
	query := "SELECT id, userid, username, channel, chat FROM chat_information" // dev:finder+query
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
			db.Logger.Debug("ok map", "values", v)
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
	// Check for errors FROM iterating over rows.
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
	query := "SELECT channel FROM chat_information" // dev:finder+query
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
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on chat_information.channel", "error", err.Error())
	}
	return info, nil
}

func (db *DBX) GetChatRunningChannels(ctx context.Context, kind string) ([]types.RunningChannels, error) {
	stats := []types.RunningChannels{}
	var query string
	switch kind {
	case "stage":
		query = "SELECT s.display_text, c.channel FROM stage AS s JOIN stage_channel AS c ON s.id = c.upstream_id WHERE s.finished = false AND c.active = true" // dev:finder+query

	case "auto_play":
		query = "SELECT a.display_text, c.channel FROM auto_play AS a JOIN auto_play_channel as c ON a.id = c.upstream_id WHERE c.active = true" // dev:finder+query

	default:
		return stats, fmt.Errorf("kind %s not found", kind)
	}
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on chat_information.channel failed", "error", err.Error())
		return stats, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.RunningChannels
		if err := rows.Scan(&s.Title, &s.Channel); err != nil {
			db.Logger.Error("scan error on chat_information.channel", "error", err.Error())
		}
		s.Kind = kind
		stats = append(stats, s)
	}

	return stats, nil
}
