package pg

import (
	"context"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetStory(ctx context.Context) ([]types.Story, error) {
	story := []types.Story{}
	query := "SELECT id, title, announcement, notes, writer_id FROM story"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on story failed", "error", err.Error())
		return story, err
	}
	defer rows.Close()
	for rows.Next() {
		var s types.Story
		if err := rows.Scan(&s.ID, &s.Title, &s.Announcement, &s.Notes, &s.WriterID); err != nil {
			db.Logger.Error("scan error on story", "error", err.Error())
		}
		story = append(story, s)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on story", "error", err.Error())
	}
	return story, nil
}

// func (db *DBX) CreateStory(ctx context.Context, title, announcement, notes string, WriterID int) (int, error) {
// 	query := "INSERT INTO story(title, announcement, notes, master_id) VALUES($1, $2, $3, $4) RETURNING id"
// 	stmt, err := db.Conn.PrepareContext(ctx, query)
// 	if err != nil {
// 		db.Logger.Error("prepare insert into story failed", "error", err.Error())
// 		return -1, err
// 	}
// 	defer stmt.Close()
// 	var res int
// 	err = stmt.QueryRow(title, announcement, notes, WriterID).Scan(&res)
// 	if err != nil {
// 		db.Logger.Error("query row insert into story failed", "error", err.Error())
// 		return -1, err
// 	}
// 	return res, nil
// }

func (db *DBX) CreateStoryTx(ctx context.Context, title, announcement, notes, encodingKey string, writerID int) (int, error) {
	// TX
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on CreateStoryTx failed", "error", err.Error())
		return -1, err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// insert story
	queryStory := "INSERT INTO story(title, notes, announcement, writer_id) VALUES($1, $2, $3, $4) RETURNING id"
	stmtStory, err := db.Conn.PrepareContext(ctx, queryStory)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer stmtStory.Close()
	var storyID int
	err = tx.StmtContext(ctx, stmtStory).QueryRow(title, notes, announcement, writerID).Scan(&storyID)
	if err != nil {
		db.Logger.Error("query row insert into story failed", "error", err.Error(), "title", title, "notes", notes, "announcement", announcement, "writerID", writerID)
		return -1, err
	}
	// insert story key
	queryKey := "INSERT INTO story_keys(encoding_key, story_id) VALUES($1, $2) RETURNING id"
	stmtStoryKeys, err := db.Conn.PrepareContext(ctx, queryKey)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer stmtStoryKeys.Close()
	var encodingKeyID int
	err = tx.StmtContext(ctx, stmtStoryKeys).QueryRow(encodingKey, storyID).Scan(&encodingKeyID)
	if err != nil {
		db.Logger.Error("query row insert into story_keys failed", "error", err.Error())
		return -1, err
	}
	// grant access to Writer to story_key
	queryAccess := "INSERT INTO access_story_keys(writer_id, story_keys_id) VALUES($1, $2) RETURNING id"
	stmtAccessStoryKeys, err := db.Conn.PrepareContext(ctx, queryAccess)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer stmtAccessStoryKeys.Close()
	var accessStoryID int
	err = tx.StmtContext(ctx, stmtAccessStoryKeys).QueryRow(writerID, encodingKeyID).Scan(&accessStoryID)
	if err != nil {
		db.Logger.Error("query row insert into access_story_keys failed", "error", err.Error())
		return -1, err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on CreateStoryTx failed", "error", err.Error())
		return -1, err
	}

	return storyID, nil
}

func (db *DBX) GetStoryIDByTitle(ctx context.Context, title string) (int, error) {
	var storyID int
	rows, err := db.Conn.QueryContext(ctx, "SELECT id FROM story WHERE title = $1", title)
	if err != nil {
		db.Logger.Error("query on story by title failed", "error", err.Error())
		return storyID, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&storyID); err != nil {
			db.Logger.Error("scan error on story by title", "error", err.Error())
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on sotry by title", "error", err.Error())
	}
	return storyID, nil
}

func (db *DBX) GetStoryByID(ctx context.Context, id int) (types.Story, error) {
	var story types.Story
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, announcement, notes, writer_id FROM story WHERE id = $1", id)
	if err != nil {
		db.Logger.Error("query on story by id failed", "error", err.Error())
		return story, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&story.ID, &story.Title, &story.Announcement, &story.Notes, &story.WriterID); err != nil {
			db.Logger.Error("scan error on story by id", "error", err.Error())
		}
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on story by id", "error", err.Error())
	}
	return story, nil
}

func (db *DBX) GetStoriesByWriterID(ctx context.Context, id int) ([]types.Story, error) {
	var stories []types.Story
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, announcement, notes, writer_id FROM story WHERE writer_id = $1", id)
	if err != nil {
		db.Logger.Error("query on story by writer_id failed", "error", err.Error())
		return stories, err
	}
	defer rows.Close()
	for rows.Next() {
		var story types.Story
		if err := rows.Scan(&story.ID, &story.Title, &story.Announcement, &story.Notes, &story.WriterID); err != nil {
			db.Logger.Error("scan error on story by writer_id", "error", err.Error())
		}
		stories = append(stories, story)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on story by Writer_id", "error", err.Error())
	}
	return stories, nil
}
