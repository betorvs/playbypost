package pg

import (
	"context"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (db *DBX) GetAllSessions(ctx context.Context) ([]types.Session, error) {
	getAllSessions := `SELECT id, username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity FROM writers_sessions` // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, getAllSessions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []types.Session
	for rows.Next() {
		var session types.Session
		err := rows.Scan(&session.ID, &session.Username, &session.Token, &session.UserID, &session.Expiry, &session.ClientType, &session.ClientInfo, &session.IPAddress, &session.UserAgent, &session.CreatedAt, &session.UpdatedAt, &session.LastActivity)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetSessionByID returns a session by its ID
func (db *DBX) GetSessionByID(ctx context.Context, sessionID int64) (types.Session, error) {
	getSessionByID := `SELECT id, username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity FROM writers_sessions WHERE id = $1` // dev:finder+query
	var session types.Session
	err := db.Conn.QueryRowContext(ctx, getSessionByID, sessionID).Scan(&session.ID, &session.Username, &session.Token, &session.UserID, &session.Expiry, &session.ClientType, &session.ClientInfo, &session.IPAddress, &session.UserAgent, &session.CreatedAt, &session.UpdatedAt, &session.LastActivity)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (db *DBX) CreateSession(ctx context.Context, session types.Session) error {
	createSession := `INSERT INTO writers_sessions (username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id` // dev:finder+query
	var sessionID int64
	err := db.Conn.QueryRowContext(ctx, createSession, session.Username, session.Token, session.UserID, session.Expiry, session.ClientType, session.ClientInfo, session.IPAddress, session.UserAgent, session.CreatedAt, session.UpdatedAt, session.LastActivity).Scan(&sessionID)
	if err != nil {
		db.Logger.Error("failed to create session", "error", err)
		return err
	}

	// Log session created event
	err = db.LogSessionCreated(ctx, session, sessionID)
	if err != nil {
		db.Logger.Error("failed to log session created event", "error", err)
		// Don't fail the session creation if event logging fails
	}

	return nil
}

func (db *DBX) GetSessionByToken(ctx context.Context, token string) (types.Session, error) {
	var session types.Session
	getSessionByToken := `SELECT id, username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity FROM writers_sessions WHERE token = $1` // dev:finder+query
	err := db.Conn.QueryRowContext(ctx, getSessionByToken, token).Scan(&session.ID, &session.Username, &session.Token, &session.UserID, &session.Expiry, &session.ClientType, &session.ClientInfo, &session.IPAddress, &session.UserAgent, &session.CreatedAt, &session.UpdatedAt, &session.LastActivity)
	if err != nil {
		return session, err
	}

	// Normalize timezone - ensure expiry time is in local timezone
	// This fixes the issue where PostgreSQL returns UTC but Go interprets as local
	if !session.Expiry.IsZero() {
		// Convert to local timezone if it's not already
		session.Expiry = session.Expiry.Local()
	}
	if !session.CreatedAt.IsZero() {
		session.CreatedAt = session.CreatedAt.Local()
	}
	if !session.UpdatedAt.IsZero() {
		session.UpdatedAt = session.UpdatedAt.Local()
	}
	if !session.LastActivity.IsZero() {
		session.LastActivity = session.LastActivity.Local()
	}

	return session, nil
}

func (db *DBX) UpdateSessionLastActivity(ctx context.Context, token string) error {
	// Get session ID and username for event logging
	sessionID, err := db.getSessionIDByToken(ctx, token)
	if err != nil {
		db.Logger.Error("failed to get session ID for activity update", "error", err)
		// Continue with update even if we can't get the session ID
	}

	updateLastActivity := `UPDATE writers_sessions SET last_activity = $1, updated_at = $1 WHERE token = $2` // dev:finder+query
	_, err = db.Conn.ExecContext(ctx, updateLastActivity, time.Now(), token)
	if err != nil {
		return err
	}

	// Log activity updated event if we have the session ID
	if sessionID > 0 {
		// Get username for event logging
		session, err := db.GetSessionByToken(ctx, token)
		if err != nil {
			db.Logger.Error("failed to get session for activity update event", "error", err)
			// Continue even if we can't get the session details
		} else {
			err = db.LogActivityUpdated(ctx, sessionID, session.Username)
			if err != nil {
				db.Logger.Error("failed to log activity updated event", "error", err)
				// Don't fail the update if event logging fails
			}
		}
	}

	return nil
}

func (db *DBX) DeleteSessionByToken(ctx context.Context, token string) error {
	// Get session ID before deletion for event logging
	sessionID, err := db.getSessionIDByToken(ctx, token)
	if err != nil {
		db.Logger.Error("failed to get session ID for deletion", "error", err)
		// Continue with deletion even if we can't get the session ID
	}

	deleteSessionByToken := `DELETE FROM writers_sessions WHERE token = $1` // dev:finder+query
	_, err = db.Conn.ExecContext(ctx, deleteSessionByToken, token)
	if err != nil {
		return err
	}

	// Log session deleted event if we have the session ID
	if sessionID > 0 {
		err = db.LogSessionDeleted(ctx, sessionID, "manual_deletion")
		if err != nil {
			db.Logger.Error("failed to log session deleted event", "error", err)
			// Don't fail the deletion if event logging fails
		}
	}

	return nil
}

func (db *DBX) DeleteSessionByID(ctx context.Context, sessionID int64) error {
	// Get session details before deletion for event logging
	var username string
	getSessionQuery := `SELECT username FROM writers_sessions WHERE id = $1` // dev:finder+query
	err := db.Conn.QueryRowContext(ctx, getSessionQuery, sessionID).Scan(&username)
	if err != nil {
		db.Logger.Error("failed to get session details for deletion", "error", err)
		// Continue with deletion even if we can't get the session details
	}

	deleteSessionByID := `DELETE FROM writers_sessions WHERE id = $1` // dev:finder+query
	_, err = db.Conn.ExecContext(ctx, deleteSessionByID, sessionID)
	if err != nil {
		return err
	}

	// Log session deleted event
	err = db.LogSessionDeleted(ctx, sessionID, "admin_deletion")
	if err != nil {
		db.Logger.Error("failed to log session deleted event", "error", err)
		// Don't fail the deletion if event logging fails
	}

	return nil
}

func (db *DBX) DeleteExpiredSessions(ctx context.Context) error {
	// First, get the count of expired sessions for logging
	var count int
	countQuery := `SELECT COUNT(*) FROM writers_sessions WHERE expiry < $1`
	err := db.Conn.QueryRowContext(ctx, countQuery, time.Now()).Scan(&count)
	if err != nil {
		db.Logger.Error("failed to count expired sessions", "error", err)
		// Continue with deletion even if counting fails
	}

	deleteExpiredSessions := `DELETE FROM writers_sessions WHERE expiry < $1` // dev:finder+query
	result, err := db.Conn.ExecContext(ctx, deleteExpiredSessions, time.Now())
	if err != nil {
		return err
	}

	// Get the actual number of rows deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		db.Logger.Error("failed to get rows affected", "error", err)
		rowsAffected = int64(count) // Fallback to the count we got earlier
	}

	// Log cleanup executed event
	err = db.LogCleanupExecuted(ctx, int(rowsAffected))
	if err != nil {
		db.Logger.Error("failed to log cleanup executed event", "error", err)
		// Don't fail the cleanup if event logging fails
	}

	return nil
}

func (db *DBX) CreateSessionEvent(ctx context.Context, event types.SessionEvent) error {
	createSessionEvent := `INSERT INTO session_events (session_id, event_type, timestamp, event_data) VALUES ($1, $2, $3, $4)` // dev:finder+query
	_, err := db.Conn.ExecContext(ctx, createSessionEvent, event.SessionID, event.EventType, event.Timestamp, event.Data)
	return err
}

func (db *DBX) GetSessionEvents(ctx context.Context) ([]types.SessionEvent, error) {
	getSessionEvents := `SELECT id, session_id, event_type, timestamp, event_data FROM session_events` // dev:finder+query
	rows, err := db.Conn.QueryContext(ctx, getSessionEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []types.SessionEvent{}
	for rows.Next() {
		var event types.SessionEvent
		err := rows.Scan(&event.ID, &event.SessionID, &event.EventType, &event.Timestamp, &event.Data)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (db *DBX) getSessionIDByToken(ctx context.Context, token string) (int64, error) {
	var sessionID int64
	getSessionIDByToken := `SELECT id FROM writers_sessions WHERE token = $1` // dev:finder+query
	err := db.Conn.QueryRowContext(ctx, getSessionIDByToken, token).Scan(&sessionID)
	if err != nil {
		return 0, err
	}
	return sessionID, nil
}
