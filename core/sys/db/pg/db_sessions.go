package pg

import (
	"context"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (c *DBX) CreateSession(ctx context.Context, session types.Session) error {
	createSession := `INSERT INTO writers_sessions (username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)` // dev:finder+query
	_, err := c.Conn.ExecContext(ctx, createSession, session.Username, session.Token, session.UserID, session.Expiry, session.ClientType, session.ClientInfo, session.IPAddress, session.UserAgent, session.CreatedAt, session.UpdatedAt, session.LastActivity)
	if err != nil {
		c.Logger.Error("failed to create session", "error", err)
		return err
	}

	return nil
}

func (c *DBX) GetSessionByToken(ctx context.Context, token string) (types.Session, error) {
	var session types.Session
	getSessionByToken := `SELECT username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity FROM writers_sessions WHERE token = $1` // dev:finder+query
	err := c.Conn.QueryRowContext(ctx, getSessionByToken, token).Scan(&session.Username, &session.Token, &session.UserID, &session.Expiry, &session.ClientType, &session.ClientInfo, &session.IPAddress, &session.UserAgent, &session.CreatedAt, &session.UpdatedAt, &session.LastActivity)
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

func (c *DBX) DeleteSessionByToken(ctx context.Context, token string) error {
	deleteSessionByToken := `DELETE FROM writers_sessions WHERE token = $1` // dev:finder+query
	_, err := c.Conn.ExecContext(ctx, deleteSessionByToken, token)
	return err
}

func (c *DBX) DeleteExpiredSessions(ctx context.Context) error {
	deleteExpiredSessions := `DELETE FROM writers_sessions WHERE expiry < $1` // dev:finder+query
	_, err := c.Conn.ExecContext(ctx, deleteExpiredSessions, time.Now())
	return err
}

func (c *DBX) CreateSessionEvent(ctx context.Context, event types.SessionEvent) error {
	createSessionEvent := `INSERT INTO session_events (session_id, event_type, timestamp, event_data) VALUES ($1, $2, $3, $4)` // dev:finder+query
	_, err := c.Conn.ExecContext(ctx, createSessionEvent, event.SessionID, event.EventType, event.Timestamp, event.Data)
	return err
}

func (c *DBX) GetSessionEvents(ctx context.Context) ([]types.SessionEvent, error) {
	getSessionEvents := `SELECT id, session_id, event_type, timestamp, event_data FROM session_events` // dev:finder+query
	rows, err := c.Conn.QueryContext(ctx, getSessionEvents)
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

func (c *DBX) GetAllSessions(ctx context.Context) ([]types.Session, error) {
	getAllSessions := `SELECT username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity FROM writers_sessions` // dev:finder+query
	rows, err := c.Conn.QueryContext(ctx, getAllSessions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []types.Session
	for rows.Next() {
		var session types.Session
		err := rows.Scan(&session.Username, &session.Token, &session.UserID, &session.Expiry, &session.ClientType, &session.ClientInfo, &session.IPAddress, &session.UserAgent, &session.CreatedAt, &session.UpdatedAt, &session.LastActivity)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
