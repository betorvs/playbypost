package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

// CreateWriterUserAssociation creates a new association between a writer and a user
func (db *DBX) CreateWriterUserAssociation(ctx context.Context, writerID, userID int) (int, error) {
	// Validate writer and user existence
	writerExists, err := db.CheckWriterExists(ctx, writerID)
	if err != nil {
		return 0, err
	}
	if !writerExists {
		return 0, fmt.Errorf("writer not found")
	}

	userExists, err := db.CheckUserExists(ctx, userID)
	if err != nil {
		return 0, err
	}
	if !userExists {
		return 0, fmt.Errorf("user not found")
	}

	// Check for duplicate association
	associationExists, err := db.CheckWriterUserAssociationExists(ctx, writerID, userID)
	if err != nil {
		return 0, err
	}
	if associationExists {
		return 0, fmt.Errorf("writer user association already exists")
	}

	query := `INSERT INTO writer_users (writer_id, user_id) VALUES ($1, $2) RETURNING id;` // dev:finder+query

	var id int
	err = db.Conn.QueryRowContext(ctx, query, writerID, userID).Scan(&id)
	if err != nil {
		db.Logger.Error(fmt.Sprintf("CreateWriterUserAssociation: %v", err))
		return 0, err
	}

	return id, nil
}

// GetUsersByWriterID retrieves all users associated with a given writer ID
func (db *DBX) GetUsersByWriterID(ctx context.Context, writerID int) ([]types.User, error) {
	// dev:finder+multiline+query
	query := `
		SELECT
			u.id, u.userid
		FROM
			users u
		JOIN
			writer_users wu ON u.id = wu.user_id
		WHERE
			wu.writer_id = $1;
	` // dev:finder+query

	rows, err := db.Conn.QueryContext(ctx, query, writerID)
	if err != nil {
		db.Logger.Error(fmt.Sprintf("GetUsersByWriterID: %v", err))
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error(fmt.Sprintf("GetUsersByWriterID rows close: %v", err))
		}
	}()

	var users []types.User
	for rows.Next() {
		var user types.User
		if err := rows.Scan(&user.ID, &user.UserID); err != nil {
			db.Logger.Error(fmt.Sprintf("GetUsersByWriterID scan: %v", err))
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		db.Logger.Error(fmt.Sprintf("GetUsersByWriterID rows error: %v", err))
		return nil, err
	}

	return users, nil
}

// GetWritersByUserID retrieves all writers associated with a given user ID
func (db *DBX) GetWritersByUserID(ctx context.Context, userID int) ([]types.Writer, error) {
	// dev:finder+multiline+query
	query := `
		SELECT
			w.id, w.username, w.password
		FROM
			writers w
		JOIN
			writer_users wu ON w.id = wu.writer_id
		WHERE
			wu.user_id = $1;
	` // dev:finder+query

	rows, err := db.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		db.Logger.Error(fmt.Sprintf("GetWritersByUserID: %v", err))
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error(fmt.Sprintf("GetWritersByUserID rows close: %v", err))
		}
	}()

	var writers []types.Writer
	for rows.Next() {
		var writer types.Writer
		if err := rows.Scan(&writer.ID, &writer.Username, &writer.Password); err != nil {
			db.Logger.Error(fmt.Sprintf("GetWritersByUserID scan: %v", err))
			return nil, err
		}
		writers = append(writers, writer)
	}

	if err := rows.Err(); err != nil {
		db.Logger.Error(fmt.Sprintf("GetWritersByUserID rows error: %v", err))
		return nil, err
	}

	return writers, nil
}

// CheckWriterExists checks if a writer with the given ID exists
func (db *DBX) CheckWriterExists(ctx context.Context, writerID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM writers WHERE id = $1);` // dev:finder+query
	var exists bool
	err := db.Conn.QueryRowContext(ctx, query, writerID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		db.Logger.Error(fmt.Sprintf("CheckWriterExists: %v", err))
		return false, err
	}
	return exists, nil
}

// CheckUserExists checks if a user with the given ID exists
func (db *DBX) CheckUserExists(ctx context.Context, userID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);` // dev:finder+query
	var exists bool
	err := db.Conn.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		db.Logger.Error(fmt.Sprintf("CheckUserExists: %v", err))
		return false, err
	}
	return exists, nil
}

// CheckWriterUserAssociationExists checks if an association between the given writer and user already exists
func (db *DBX) CheckWriterUserAssociationExists(ctx context.Context, writerID, userID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM writer_users WHERE writer_id = $1 AND user_id = $2);` // dev:finder+query
	var exists bool
	err := db.Conn.QueryRowContext(ctx, query, writerID, userID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		db.Logger.Error(fmt.Sprintf("CheckWriterUserAssociationExists: %v", err))
		return false, err
	}
	return exists, nil
}

// GetWriterUsersAssociation retrieves all associations between a writer and a user
func (db *DBX) GetWriterUsersAssociation(ctx context.Context) ([]types.WriterUserAssociation, error) {
	query := `SELECT id, writer_id, user_id FROM writer_users;` // dev:finder+query
	var associations []types.WriterUserAssociation
	rows, err := db.Conn.QueryContext(ctx, query)
	if err != nil {
		db.Logger.Error(fmt.Sprintf("GetWriterUsersAssociation: %v", err))
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			db.Logger.Error(fmt.Sprintf("GetWriterUsersAssociation rows close: %v", err))
		}
	}()
	for rows.Next() {
		var association types.WriterUserAssociation
		if err := rows.Scan(&association.ID, &association.WriterID, &association.UserID); err != nil {
			db.Logger.Error(fmt.Sprintf("GetWriterUsersAssociation scan: %v", err))
			return nil, err
		}
		db.Logger.Info("association", "id", association.ID, "writer_id", association.WriterID, "user_id", association.UserID)
		if association.ID != 0 {
			associations = append(associations, association)
		}
	}
	if err := rows.Err(); err != nil {
		db.Logger.Error(fmt.Sprintf("GetWriterUsersAssociation rows error: %v", err))
		return nil, err
	}
	return associations, nil
}

// DeleteWriterUserAssociation deletes an association between a writer and a user
func (db *DBX) DeleteWriterUserAssociation(ctx context.Context, id int) error {
	query := `DELETE FROM writer_users WHERE id = $1;` // dev:finder+query
	_, err := db.Conn.ExecContext(ctx, query, id)
	if err != nil {
		db.Logger.Error(fmt.Sprintf("DeleteWriterUserAssociation: %v", err))
		return err
	}
	return nil
}
