package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) CreateTask(ctx context.Context, description, ability, skill string, kind types.TaskKind, target int) (int, error) {
	query := "INSERT INTO tasks(description, kind, ability, skill, target) VALUES($1, $2, $3, $4, $5) RETURNING id" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare insert into tasks failed", "error", err.Error())
		return -1, err
	}
	defer stmt.Close()
	var res int
	err = stmt.QueryRow(description, kind, ability, skill, target).Scan(&res)
	if err != nil {
		db.Logger.Error("query row insert into tasks failed", "error", err.Error())
		return -1, db.parsePostgresError(err)
	}
	return res, nil
}

func (db *DBX) GetTask(ctx context.Context) ([]types.Task, error) {
	t := []types.Task{}
	query := "SELECT id, description, kind, ability, skill, target FROM tasks" // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error("query on tasks failed", "error", err.Error())
		return t, err
	}
	defer rows.Close()
	for rows.Next() {
		var tl types.Task
		if err := rows.Scan(&tl.ID, &tl.Description, &tl.Kind, &tl.Ability, &tl.Skill, &tl.Target); err != nil {
			db.Logger.Error("scan error on tasks", "error", err.Error())
		}
		t = append(t, tl)
	}
	// Check for errors FROM iterating over rows.
	if err := rows.Err(); err != nil {
		db.Logger.Error("rows err on tasks", "error", err.Error())
	}
	return t, nil
}

func (db *DBX) GetTaskByID(ctx context.Context, id int) (types.Task, error) {
	t := types.Task{}
	query := "SELECT id, description, kind, ability, skill, target FROM tasks WHERE id = $1" // dev:finder+query
	err := db.Conn.QueryRowContext(ctx, query, id).Scan(&t.ID, &t.Description, &t.Kind, &t.Ability, &t.Skill, &t.Target)
	if err != nil {
		db.Logger.Error("query on tasks by id failed", "error", err.Error())
		return t, err
	}
	return t, nil
}

func (db *DBX) UpdateTaskByID(ctx context.Context, description, ability, skill string, kind types.TaskKind, target, id int) error {
	query := "UPDATE tasks SET description = $1, ability = $2, skill = $3, kind = $4, target = $5 WHERE id = $6" // dev:finder+query
	stmt, err := db.Conn.PrepareContext(ctx, query)
	if err != nil {
		db.Logger.Error("prepare update tasks failed", "error", err.Error())
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, description, ability, skill, kind, target, id)
	if err != nil {
		db.Logger.Error("exec update tasks failed", "error", err.Error())
		return err
	}
	return nil
}

func (db *DBX) DeleteTaskByID(ctx context.Context, id int) error {
	// tx
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		db.Logger.Error("tx begin on DeleteTaskByID failed", "error", err.Error())
		return err
	}
	// Defer a rollback in case anything fails.
	defer func() {
		rollback := tx.Rollback()
		if err != nil && rollback != nil {
			err = fmt.Errorf("rolling back transaction: %w", err)
		}
	}()
	// check if this have an stage_running_tasks
	var countStageRunning int
	queryCheck := "SELECT COUNT(*) FROM stage_running_tasks WHERE task_id = $1" // dev:finder+query
	if err = tx.QueryRowContext(ctx, queryCheck, id).Scan(&countStageRunning); err != nil {
		if err != sql.ErrNoRows {
			db.Logger.Error("no rows passed", "err", err.Error())
			return err
		}
	}
	if countStageRunning > 0 {
		db.Logger.Error("found task", "task_id", id)
		return fmt.Errorf("found tasks associated with stage in stage_running_tasks")
	}
	// delete task
	queryTask := "DELETE FROM tasks WHERE id = $1" // dev:finder+query
	stmtTask, err := db.Conn.PrepareContext(ctx, queryTask)
	if err != nil {
		db.Logger.Error("tx prepare on tasks failed", "error", err.Error())
		return err
	}
	defer stmtTask.Close()
	_, err = tx.StmtContext(ctx, stmtTask).ExecContext(ctx, id)
	if err != nil {
		db.Logger.Error("exec on tasks failed", "error", err.Error())
		return err
	}
	// commit if everything is okay
	if err = tx.Commit(); err != nil {
		db.Logger.Error("tx commit on DeleteTaskByID failed", "error", err.Error())
		return err
	}
	return nil
}
