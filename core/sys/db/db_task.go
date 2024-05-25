package db

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) CreateTask(ctx context.Context, title, displayText, checks string, kind, target, encounterID int, options map[string]int) (int, error) {
	query := "INSERT INTO tasks(title, encounters_id, display_text, kind, checks, target, options, finished) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("prepare insert into tasks failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(title, encounterID, displayText, kind, checks, target, options, false).Scan(&res)
	if err != nil {
		db.logger.Error("query row insert into tasks failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) GetTasksByEncounterID(ctx context.Context, id int) (map[string]types.Task, error) {
	tasks := make(map[string]types.Task)
	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, display_text, kind, checks, target, options, finished FROM tasks WHERE encounters_id = $1", id)
	if err != nil {
		db.logger.Error("query on tasks by encounter id failed", "error", err.Error())
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var task types.Task
		if err := rows.Scan(&id, &task.Title, &task.DisplayText, &task.Kind, &task.Checks, &task.Target, &task.Options, &task.Finished); err != nil {
			db.logger.Error("scan error on tasks by id", "error", err.Error())
		}
		tasks[task.DisplayText] = task
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows error on story by id", "error", err.Error())
	}
	return tasks, nil
}
