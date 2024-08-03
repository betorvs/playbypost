package db

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) CreateTask(ctx context.Context, description, ability, skill string, kind types.TaskKind, target int) (int, error) {
	query := "INSERT INTO tasks(description, kind, ability, skill, target) VALUES($1, $2, $3, $4, $5) RETURNING id"
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.logger.Error("prepare insert into tasks failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(description, kind, ability, skill, target).Scan(&res)
	if err != nil {
		db.logger.Error("query row insert into tasks failed", "error", err.Error())
		return -1, err
	}
	return res, nil
}

func (db *DBX) GetTask(ctx context.Context) ([]types.Task, error) {
	t := []types.Task{}
	query := "SELECT id, description, kind, ability, skill, target FROM tasks"
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.logger.Error("query on tasks failed", "error", err.Error())
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		var tl types.Task
		if err := rows.Scan(&tl.ID, &tl.Description, &tl.Kind, &tl.Ability, &tl.Skill, &tl.Target); err != nil {
			db.logger.Error("scan error on tasks", "error", err.Error())
		}
		t = append(t, tl)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		db.logger.Error("rows err on tasks", "error", err.Error())
	}
	return t, nil
}

// func (db *DBX) GetTasksByEncounterID(ctx context.Context, id int) (map[string]types.Task, error) {
// 	tasks := make(map[string]types.Task)
// 	rows, err := db.Conn.QueryContext(ctx, "SELECT id, title, display_text, kind, checks, target, options, finished FROM tasks WHERE encounters_id = $1", id)
// 	if err != nil {
// 		db.logger.Error("query on tasks by encounter id failed", "error", err.Error())
// 		return tasks, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int
// 		var task types.Task
// 		if err := rows.Scan(&id, &task.Title, &task.DisplayText, &task.Kind, &task.Checks, &task.Target, &task.Options, &task.Finished); err != nil {
// 			db.logger.Error("scan error on tasks by id", "error", err.Error())
// 		}
// 		tasks[task.DisplayText] = task
// 	}
// 	// Check for errors from iterating over rows.
// 	if err := rows.Err(); err != nil {
// 		db.logger.Error("rows error on story by id", "error", err.Error())
// 	}
// 	return tasks, nil
// }
