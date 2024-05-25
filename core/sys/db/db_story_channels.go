package db

import (
	"context"
)

func (db *DBX) GetStoryChannels(ctx context.Context) (map[string]int, error) {
	channels := make(map[string]int)
	query := "SELECT c.id, c.story_id, c.channel FROM story_channels AS c JOIN story AS s ON s.id = c.story_id WHERE s.finished = false"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.logger.Error("query on story channels failed", "error", err.Error())
		return channels, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, storyID int
		var channel string
		if err := rows.Scan(&id, &storyID, &channel); err != nil {
			db.logger.Error("scan error on story channels", "error", err.Error())
		}
		value, ok := channels[channel]
		if !ok {
			channels[channel] = storyID
		} else {
			db.logger.Error("channel added to story duplicate", "channel_id", value, "id", id)
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows err on story", "error", err.Error())
	}
	return channels, nil
}
