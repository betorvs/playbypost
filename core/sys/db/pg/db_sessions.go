package pg

import (
	"context"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

func (c *DBX) CreateSession(ctx context.Context, session types.Session) error {
	createSession := `INSERT INTO writers_sessions (username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)` // dev:finder+query
	_, err := c.Conn.ExecContext(ctx, createSession, session.Username, session.Token, session.UserID, session.Expiry, session.ClientType, session.ClientInfo, session.IPAddress, session.UserAgent, session.CreatedAt, session.UpdatedAt, session.LastActivity)
	return err
}

func (c *DBX) GetSessionByToken(ctx context.Context, token string) (types.Session, error) {
	var session types.Session
	getSessionByToken := `SELECT id, username, token, user_id, expiry, client_type, client_info, ip_address, user_agent, created_at, updated_at, last_activity FROM writers_sessions WHERE token = $1` // dev:finder+query
	err := c.Conn.QueryRowContext(ctx, getSessionByToken, token).Scan(&session.Username, &session.Token, &session.UserID, &session.Expiry, &session.ClientType, &session.ClientInfo, &session.IPAddress, &session.UserAgent, &session.CreatedAt, &session.UpdatedAt, &session.LastActivity)
	return session, err
}

func (c *DBX) DeleteSessionByToken(ctx context.Context, token string) error {
	deleteSessionByToken := `DELETE FROM writers_sessions WHERE token = $1` // dev:finder+query
	_, err := c.Conn.ExecContext(ctx, deleteSessionByToken, token)
	return err
}
