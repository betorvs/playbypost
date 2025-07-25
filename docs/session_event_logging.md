# Session Event Logging and Maintenance Tools

## Purpose
This document outlines the comprehensive session event logging system and maintenance tools implemented in the Play-by-Post application. The system provides detailed tracking of session lifecycle events, enabling administrators to monitor user activity, troubleshoot authentication issues, and perform maintenance operations.

## Architecture Overview

### Event Logging Integration
The session event logging system is integrated directly into the database layer to avoid import loops and maintain clean separation of concerns. Events are logged within the same transaction as the corresponding data operations, ensuring consistency and reliability.

### Key Components
- **Database Layer Integration**: Event logging is embedded directly in PostgreSQL operations
- **Structured Event Data**: JSON-formatted event data with relevant context
- **Non-blocking Logging**: Event logging failures don't affect primary operations
- **Maintenance Tools**: CLI and API tools for session management

## Session Events Table Structure

The `session_events` table stores comprehensive event data with the following structure:

- `id`: bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY. A unique identifier for each event.
- `session_id`: BIGINT NOT NULL REFERENCES writers_sessions(id) ON DELETE CASCADE. The ID of the associated session.
- `event_type`: VARCHAR(50) NOT NULL. The type of event (e.g., 'session_created', 'session_deleted').
- `event_data`: JSONB. Structured JSON data containing event-specific information.
- `timestamp`: TIMESTAMP NOT NULL DEFAULT NOW(). The timestamp when the event occurred.

## Event Types

The system tracks the following session lifecycle events:

### Session Creation Events
- **`session_created`**: Logged when a new session is created
  - Data includes: username, client_type, ip_address, user_agent

### Session Deletion Events
- **`session_deleted`**: Logged when a session is manually deleted
  - Data includes: username, deletion_reason

### Session Expiration Events
- **`session_expired`**: Logged when sessions are automatically expired
  - Data includes: count of expired sessions, cleanup_timestamp

### Authentication Events
- **`login_attempt`**: Logged for all login attempts (success and failure)
  - Data includes: username, success_status, ip_address, user_agent
- **`logout`**: Logged when users explicitly log out
  - Data includes: username, session_duration

### Session Validation Events
- **`session_validated`**: Logged when sessions are validated
  - Data includes: validation_source (cache/database), username
- **`session_invalid`**: Logged when session validation fails
  - Data includes: failure_reason (not_found/expired)

### Activity Tracking Events
- **`activity_updated`**: Logged when session activity is updated
  - Data includes: username, last_activity_timestamp

### Maintenance Events
- **`cleanup_executed`**: Logged when session cleanup operations are performed
  - Data includes: sessions_removed_count, cleanup_timestamp

## Database Interface

The `DBClient` interface includes the following methods for session event logging:

```go
type DBClient interface {
    // Session operations with integrated event logging
    CreateSession(ctx context.Context, session types.Session) error
    DeleteSessionByToken(ctx context.Context, token string) error
    DeleteSessionByID(ctx context.Context, sessionID int64) error
    DeleteExpiredSessions(ctx context.Context) ([]types.Session, error)
    UpdateSessionLastActivity(ctx context.Context, token string) error
    
    // Event logging operations
    CreateSessionEvent(ctx context.Context, event types.SessionEvent) error
    GetSessionEvents(ctx context.Context) ([]types.SessionEvent, error)
}
```

## Event Data Structure

Events use structured JSON data to provide context and enable analysis:

```go
type SessionEvent struct {
    ID        int       `json:"id"`
    SessionID int64     `json:"session_id"`
    EventType string    `json:"event_type"`
    Timestamp time.Time `json:"timestamp"`
    Data      string    `json:"data"` // JSON string
}
```

### Example Event Data

```json
{
  "username": "storyteller123",
  "client_type": "browser",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "session_duration": "2h30m"
}
```

## Maintenance Tools

### CLI Commands (admin-ctl)

The `admin-ctl` tool provides comprehensive session management capabilities:

#### List All Sessions
```bash
./admin-ctl session list
```
Displays all active sessions with details including username, client type, creation time, and last activity.

#### Cleanup Expired Sessions
```bash
./admin-ctl session cleanup
```
Manually triggers the cleanup of expired sessions and displays the count of removed sessions.

#### View Session Events
```bash
./admin-ctl session events
```
Lists all session events with timestamps, event types, and associated data.

#### Show Session Statistics
```bash
./admin-ctl session stats
```
Displays session statistics including total sessions, active sessions, and recent activity.

#### Delete Specific Session
```bash
./admin-ctl session delete --session-id 123
```
Deletes a specific session by ID and logs the deletion event.

### API Endpoints

The system provides RESTful API endpoints for session management:

#### GET /api/v1/session
Retrieves all active sessions.

#### GET /api/v1/session/events
Retrieves all session events.

#### GET /api/v1/session/active
Retrieves currently active sessions.

#### PUT /api/v1/session/cleanup
Manually triggers session cleanup and returns the count of removed sessions.

#### DELETE /api/v1/session/{id}
Deletes a specific session by ID.

## Error Handling

### Non-blocking Event Logging
Event logging failures do not affect primary operations:

```go
// Log event but don't fail the operation if logging fails
err = c.CreateSessionEvent(ctx, event)
if err != nil {
    c.Logger.Error("failed to log session event", "error", err, "event_type", event.EventType)
    // Continue with primary operation
}
```

### Graceful Degradation
- If event logging fails, the primary operation continues
- Event logging errors are logged for debugging
- System maintains functionality even with logging issues

## Migration Files

The session events table is created using the following migration:

- `core/sys/db/data/migrations/000005_create_session_events_table.up.sql`: Creates the session_events table
- `core/sys/db/data/migrations/000005_create_session_events_table.down.sql`: Drops the session_events table

### Running the Migration

```bash
# Apply the migration
./admin-ctl db up

# Rollback the migration
./admin-ctl db down
```

## Usage Examples

### Monitoring User Activity
```bash
# Check recent login attempts
./admin-ctl session events | grep "login_attempt"

# View session statistics
./admin-ctl session stats

# Monitor active sessions
./admin-ctl session list
```

### Maintenance Operations
```bash
# Clean up expired sessions
./admin-ctl session cleanup

# Remove a specific problematic session
./admin-ctl session delete --session-id 456

# Review session events for troubleshooting
./admin-ctl session events | grep "session_invalid"
```

### API Usage
```bash
# Get all sessions via API
curl -X GET http://localhost:8080/api/v1/session

# Trigger cleanup via API
curl -X PUT http://localhost:8080/api/v1/session/cleanup

# Delete a session via API
curl -X DELETE http://localhost:8080/api/v1/session/123
```

## Best Practices

### Event Data Design
- Keep event data concise but informative
- Use consistent field names across event types
- Include relevant context for debugging and analysis
- Avoid sensitive information in event data

### Performance Considerations
- Event logging is designed to be lightweight
- Events are logged asynchronously to avoid blocking
- Database indexes optimize event queries
- Regular cleanup prevents table bloat

### Security Considerations
- Event data is logged but not exposed to end users
- Sensitive information is excluded from event data
- Access to maintenance tools requires proper authentication
- Event logs are retained for debugging and compliance

## Troubleshooting

### Common Issues

#### Event Logging Failures
- Check database connectivity
- Verify session_events table exists
- Review database logs for constraint violations
- Ensure proper JSON formatting in event data

#### Session Cleanup Issues
- Verify session expiry logic
- Check for long-running transactions
- Review cleanup frequency and timing
- Monitor cleanup performance

#### Maintenance Tool Errors
- Verify admin-ctl permissions
- Check API endpoint availability
- Review authentication requirements
- Ensure proper session ID format

### Debugging Commands
```bash
# Check database connectivity
./admin-ctl db ping

# Verify migration status
./admin-ctl db version

# Review recent events
./admin-ctl session events | tail -20

# Check session table structure
./admin-ctl db inspect writers_sessions
``` 