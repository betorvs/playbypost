package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetStory(ctx context.Context) ([]types.Story, error) {
	story := []types.Story{}
	query := "SELECT id, title, announcement, notes, writer_id FROM story" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on story failed", "error", err.Error())
		return story, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		var s types.Story
		if err := rows.Scan(&s.ID, &s.Title, &s.Announcement, &s.Notes, &s.WriterID); err != nil {
			db.Logger.Error("scan error on story", "error", err.Error())
		}
		story = append(story, s)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on story", "error", err.Error())
	}
	return story, nil
}

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
	queryStory := "INSERT INTO story(title, notes, announcement, writer_id) VALUES($1, $2, $3, $4) RETURNING id" // dev:finder+query
	stmtStory, err := db.Conn.PrepareContext(ctx, queryStory)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmtStory.Close()
		if err != nil {
			db.Logger.Error("error closing stmtStory", "error", err)
		}
	}()
	var storyID int
	err = tx.StmtContext(ctx, stmtStory).QueryRow(title, notes, announcement, writerID).Scan(&storyID)
	if err != nil {
		db.Logger.Error("query row insert into story failed", "error", err.Error(), "title", title, "notes", notes, "announcement", announcement, "writerID", writerID)
		return -1, db.parsePostgresError(err)
	}
	// insert story key
	queryKey := "INSERT INTO story_keys(encoding_key, story_id) VALUES($1, $2) RETURNING id" // dev:finder+query
	stmtStoryKeys, err := db.Conn.PrepareContext(ctx, queryKey)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmtStoryKeys.Close()
		if err != nil {
			db.Logger.Error("error closing stmtStoryKeys", "error", err)
		}
	}()
	var encodingKeyID int
	err = tx.StmtContext(ctx, stmtStoryKeys).QueryRow(encodingKey, storyID).Scan(&encodingKeyID)
	if err != nil {
		db.Logger.Error("query row insert into story_keys failed", "error", err.Error())
		return -1, db.parsePostgresError(err)
	}
	// grant access to Writer to story_key
	queryAccess := "INSERT INTO access_story_keys(writer_id, story_keys_id) VALUES($1, $2) RETURNING id" // dev:finder+query
	stmtAccessStoryKeys, err := db.Conn.PrepareContext(ctx, queryAccess)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmtAccessStoryKeys.Close()
		if err != nil {
			db.Logger.Error("error closing stmtAccessStoryKeys", "error", err)
		}
	}()
	var accessStoryID int
	err = tx.StmtContext(ctx, stmtAccessStoryKeys).QueryRow(writerID, encodingKeyID).Scan(&accessStoryID)
	if err != nil {
		db.Logger.Error("query row insert into access_story_keys failed", "error", err.Error())
		return -1, db.parsePostgresError(err)
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on CreateStoryTx failed", "error", err.Error())
		return -1, err
	}

	return storyID, nil
}

func (db *DBX) UpdateStoryTx(ctx context.Context, title, announcement, notes string, id int) (int, error) {
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
	queryStory := "UPDATE story SET title = $1, notes = $2, announcement = $3 WHERE id = $4 RETURNING id" // dev:finder+query
	stmtStory, err := db.Conn.PrepareContext(ctx, queryStory)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return -1, err
	}
	defer func() {
		err := stmtStory.Close()
		if err != nil {
			db.Logger.Error("error closing stmtStory", "error", err)
		}
	}()
	var storyID int
	err = tx.StmtContext(ctx, stmtStory).QueryRow(title, notes, announcement, id).Scan(&storyID)
	if err != nil {
		db.Logger.Error("query row insert into story failed", "error", err.Error(), "notes", notes, "announcement", announcement)
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
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(&storyID); err != nil {
			db.Logger.Error("scan error on story by title", "error", err.Error())
		}
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on sotry by title", "error", err.Error())
	}
	return storyID, nil
}

func (db *DBX) GetStoryByID(ctx context.Context, id int) (types.Story, error) {
	story := types.Story{}
	query := "SELECT id, title, announcement, notes, writer_id FROM story WHERE id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on story by id failed", "error", err.Error())
		return story, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		if err := rows.Scan(&story.ID, &story.Title, &story.Announcement, &story.Notes, &story.WriterID); err != nil {
			db.Logger.Error("scan error on story by id", "error", err.Error())
		}
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on story by id", "error", err.Error())
	}
	return story, nil
}

func (db *DBX) GetStoriesByWriterID(ctx context.Context, id int) ([]types.Story, error) {
	stories := []types.Story{}
	query := "SELECT id, title, announcement, notes, writer_id FROM story WHERE writer_id = $1" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on story by writer_id failed", "error", err.Error())
		return stories, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	for rows.Next() {
		var story types.Story
		if err := rows.Scan(&story.ID, &story.Title, &story.Announcement, &story.Notes, &story.WriterID); err != nil {
			db.Logger.Error("scan error on story by writer_id", "error", err.Error())
		}
		stories = append(stories, story)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows error on story by Writer_id", "error", err.Error())
	}
	return stories, nil
}

func (db *DBX) DeleteStoryByID(ctx context.Context, id int) error {
	// tx
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on DeleteStoryByID failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// check if this have an encounter
	var countEncounters int
	queryCheckEncounter := "SELECT COUNT(*) FROM encounters WHERE story_id = $1" // dev:finder+query
	if err = tx.QueryRowContext(ctx, queryCheckEncounter, id).Scan(&countEncounters); err != nil {
		if err != sql.ErrNoRows {
			db.Logger.Error("no rows passed", "err", err.Error())
			return err
		}
	}
	if countEncounters > 0 {
		db.Logger.Error("found encounters", "story_id", id)
		return fmt.Errorf("found encounters with this story")
	}
	// select story_keys_id
	query := "SELECT a.id, k.id FROM access_story_keys AS a JOIN story_keys AS k ON k.id = a.story_keys_id WHERE k.story_id = $1" // dev:finder+query
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		db.Logger.Error("query on access_story_keys failed", "error", err.Error())
		return err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error("error closing rows", "error", err)
		}
	}()
	var accessStoryID, storyKeyID int
	for rows.Next() {
		if err := rows.Scan(&accessStoryID, &storyKeyID); err != nil {
			db.Logger.Error("scan error on access_story_keys", "error", err.Error())
		}
	}
	// delete access_story_keys
	queryAccess := "DELETE FROM access_story_keys WHERE id = $1" // dev:finder+query
	stmtAccessStoryKeys, err := db.Conn.PrepareContext(ctx, queryAccess)
	if err != nil {
		db.Logger.Error("tx prepare on access_story_keys failed", "error", err.Error())
		return err
	}
	defer func() {
		err := stmtAccessStoryKeys.Close()
		if err != nil {
			db.Logger.Error("error closing stmtAccessStoryKeys", "error", err)
		}
	}()
	_, err = tx.StmtContext(ctx, stmtAccessStoryKeys).ExecContext(ctx, accessStoryID)
	if err != nil {
		db.Logger.Error("exec on access_story_keys failed", "error", err.Error())
		return err
	}
	// delete story_keys
	queryKey := "DELETE FROM story_keys WHERE id = $1" // dev:finder+query
	stmtStoryKeys, err := db.Conn.PrepareContext(ctx, queryKey)
	if err != nil {
		db.Logger.Error("tx prepare on story_keys failed", "error", err.Error())
		return err
	}
	defer func() {
		err := stmtStoryKeys.Close()
		if err != nil {
			db.Logger.Error("error closing stmtStoryKeys", "error", err)
		}
	}()
	_, err = tx.StmtContext(ctx, stmtStoryKeys).ExecContext(ctx, storyKeyID)
	if err != nil {
		db.Logger.Error("exec on story_keys failed", "error", err.Error())
		return err
	}
	// delete story
	queryStory := "DELETE FROM story WHERE id = $1" // dev:finder+query
	stmtStory, err := db.Conn.PrepareContext(ctx, queryStory)
	if err != nil {
		db.Logger.Error("tx prepare on story failed", "error", err.Error())
		return err
	}
	defer func() {
		err := stmtStory.Close()
		if err != nil {
			db.Logger.Error("error closing stmtStory", "error", err)
		}
	}()
	_, err = tx.StmtContext(ctx, stmtStory).ExecContext(ctx, id)
	if err != nil {
		db.Logger.Error("exec on story failed", "error", err.Error())
		return err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on DeleteStoryByID failed", "error", err.Error())
		return err
	}
	return nil

}
