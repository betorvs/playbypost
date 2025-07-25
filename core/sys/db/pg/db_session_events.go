package pg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/betorvs/playbypost/core/sys/web/types"
)

// Event types for session events
const (
	EventTypeSessionCreated   = "session_created"
	EventTypeSessionUpdated   = "session_updated"
	EventTypeSessionExpired   = "session_expired"
	EventTypeSessionDeleted   = "session_deleted"
	EventTypeLoginAttempt     = "login_attempt"
	EventTypeLoginSuccess     = "login_success"
	EventTypeLoginFailed      = "login_failed"
	EventTypeLogout           = "logout"
	EventTypeSessionValidated = "session_validated"
	EventTypeSessionInvalid   = "session_invalid"
	EventTypeActivityUpdated  = "activity_updated"
	EventTypeCleanupExecuted  = "cleanup_executed"
)

// LogSessionCreated logs a session creation event
func (db *DBX) LogSessionCreated(ctx context.Context, session types.Session, sessionID int64) error {
	eventData := map[string]interface{}{
		"username":    session.Username,
		"user_id":     session.UserID,
		"client_type": session.ClientType,
		"client_info": session.ClientInfo,
		"ip_address":  session.IPAddress,
		"user_agent":  session.UserAgent,
		"expiry":      session.Expiry,
		"created_at":  session.CreatedAt,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal session created event data", "error", err)
		return err
	}

	event := types.SessionEvent{
		SessionID: sessionID,
		EventType: EventTypeSessionCreated,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogSessionDeleted logs a session deletion event
func (db *DBX) LogSessionDeleted(ctx context.Context, sessionID int64, reason string) error {
	eventData := map[string]interface{}{
		"reason": reason,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal session deleted event data", "error", err)
		return err
	}

	event := types.SessionEvent{
		SessionID: sessionID,
		EventType: EventTypeSessionDeleted,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogSessionExpired logs a session expiration event
func (db *DBX) LogSessionExpired(ctx context.Context, sessionID int64) error {
	event := types.SessionEvent{
		SessionID: sessionID,
		EventType: EventTypeSessionExpired,
		Timestamp: time.Now(),
		Data:      "{}",
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogLoginAttempt logs a login attempt event
func (db *DBX) LogLoginAttempt(ctx context.Context, username, ipAddress, userAgent string, success bool) error {
	eventData := map[string]interface{}{
		"username":   username,
		"ip_address": ipAddress,
		"user_agent": userAgent,
		"success":    success,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal login attempt event data", "error", err)
		return err
	}

	eventType := EventTypeLoginSuccess
	if !success {
		eventType = EventTypeLoginFailed
	}

	event := types.SessionEvent{
		SessionID: 0, // No session ID for login attempts
		EventType: eventType,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogLogout logs a logout event
func (db *DBX) LogLogout(ctx context.Context, sessionID int64, username string) error {
	eventData := map[string]interface{}{
		"username": username,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal logout event data", "error", err)
		return err
	}

	event := types.SessionEvent{
		SessionID: sessionID,
		EventType: EventTypeLogout,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogSessionValidated logs a session validation event
func (db *DBX) LogSessionValidated(ctx context.Context, sessionID int64, username string) error {
	eventData := map[string]interface{}{
		"username": username,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal session validated event data", "error", err)
		return err
	}

	event := types.SessionEvent{
		SessionID: sessionID,
		EventType: EventTypeSessionValidated,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogSessionInvalid logs a session invalidation event
func (db *DBX) LogSessionInvalid(ctx context.Context, sessionID int64, reason string) error {
	eventData := map[string]interface{}{
		"reason": reason,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal session invalid event data", "error", err)
		return err
	}

	event := types.SessionEvent{
		SessionID: sessionID,
		EventType: EventTypeSessionInvalid,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogActivityUpdated logs a session activity update event
func (db *DBX) LogActivityUpdated(ctx context.Context, sessionID int64, username string) error {
	eventData := map[string]interface{}{
		"username": username,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal activity updated event data", "error", err)
		return err
	}

	event := types.SessionEvent{
		SessionID: sessionID,
		EventType: EventTypeActivityUpdated,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}

// LogCleanupExecuted logs a cleanup execution event
func (db *DBX) LogCleanupExecuted(ctx context.Context, sessionsDeleted int) error {
	eventData := map[string]interface{}{
		"sessions_deleted": sessionsDeleted,
	}

	data, err := json.Marshal(eventData)
	if err != nil {
		db.Logger.Error("failed to marshal cleanup executed event data", "error", err)
		return err
	}

	event := types.SessionEvent{
		SessionID: 0, // No specific session for cleanup events
		EventType: EventTypeCleanupExecuted,
		Timestamp: time.Now(),
		Data:      string(data),
	}

	return db.CreateSessionEvent(ctx, event)
}
