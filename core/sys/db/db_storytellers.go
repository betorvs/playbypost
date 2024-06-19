package db

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) CreateStorytellers(ctx context.Context, username, password string) (int, error) {
	query := "INSERT INTO writers(username, password) VALUES($1, $2) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("prepare insert into writers failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(username, password).Scan(&res)
	if err != nil {
		db.logger.Error("query row insert into writers failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) GetStorytellers(ctx context.Context, active bool) ([]types.Storyteller, error) {
	query := "SELECT id, username FROM writers"
	if active {
		query = "SELECT id, username FROM writers WHERE active = true"
	}
	users := []types.Storyteller{}
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.logger.Error("query on writers failed", "error", err.Error())
		return users, err
	}
	defer rows.Close()
	for rows.Next() {
		var user types.Storyteller
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			db.logger.Error("scan error on writers", "error", err.Error())
		}
		users = append(users, user)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows error on writers", "error", err.Error())
	}
	return users, nil
}

func (db *DBX) GetStorytellerByID(ctx context.Context, id int) (types.Storyteller, error) {
	user := types.Storyteller{}
	keys := make(map[int]string)
	rows, err := db.Conn.QueryContext(ctx, "SELECT w.id, w.username, k.story_id, k.encoding_key FROM writers AS w JOIN access_story_keys AS a ON w.id = a.writer_id JOIN story_keys AS k ON a.story_keys_id = k.id WHERE w.id = $1", id)
	if err != nil {
		db.logger.Error("query on writers by id failed", "error", err.Error())
		return user, err
	}
	defer rows.Close()
	for rows.Next() {
		// var user types.Storyteller
		var key string
		var keyID int
		if err := rows.Scan(&user.ID, &user.Username, &keyID, &key); err != nil {
			db.logger.Error("scan error on writers by id", "error", err.Error())
		}
		_, ok := keys[keyID]
		if !ok {
			keys[keyID] = key
		}
		// users = append(users, user)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows error on writers by id", "error", err.Error())
	}
	user.EncodingKeys = keys
	return user, nil
}

func (db *DBX) GetStorytellerByUsername(ctx context.Context, username string) (types.Storyteller, error) {
	user := types.Storyteller{}
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, username, password FROM writers WHERE username = $1", username)
	if err != nil {
		db.logger.Error("query on writers by username failed", "error", err.Error())
		return user, err
	}
	defer rows.Close()
	for rows.Next() {
		// var user types.Storyteller
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			db.logger.Error("scan error on writers by username", "error", err.Error())
		}
		// users = append(users, user)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows error on writers by username", "error", err.Error())
	}
	return user, nil
}

// func (db *DBX) GetUserCard(ctx context.Context) ([]types.Card, error) {
// 	query := "SELECT u.id, u.username, u.userid, s.title, p.character_name FROM writers AS u LEFT JOIN story AS s ON u.id = s.master_id LEFT JOIN players AS p ON u.id = p.player_id"
// 	users := []types.Card{}
// 	rows, err := db.Conn.QueryContext(ctx, query)
// 	if err != nil {
// 		db.logger.Error("query with joins on writers failed", "error", err.Error())
// 		return users, err
// 	}
// 	defer rows.Close()
// 	userCard := make(map[int]types.Card)
// 	for rows.Next() {
// 		var id int
// 		// var user types.Card
// 		var user, userID, title, player sql.NullString
// 		if err := rows.Scan(&id, &user, &userID, &title, &player); err != nil {
// 			db.logger.Error("scan error on writers with joins", "error", err.Error())
// 		}
// 		if v, ok := userCard[id]; ok {
// 			if title.Valid {
// 				if !slices.Contains(v.Stories, title.String) {
// 					v.Stories = append(v.Stories, title.String)
// 				}
// 			}
// 			if player.Valid {
// 				if !slices.Contains(v.Players, player.String) {
// 					v.Players = append(v.Players, player.String)
// 				}
// 			}
// 			userCard[id] = v
// 		} else {
// 			userCard[id] = types.Card{
// 				Username: user.String,
// 				UserID:   user.String,
// 				Stories:  []string{title.String},
// 				Players:  []string{player.String},
// 			}
// 		}

// 		// users = append(users, user)
// 	}
// 	// Check for errors from iterating over rows.
// 	if err := rows.Err(); err != nil {
// 		db.logger.Error("rows error on writers with joins", "error", err.Error())
// 	}
// 	for _, v := range userCard {
// 		users = append(users, v)
// 	}
// 	return users, nil
// }
