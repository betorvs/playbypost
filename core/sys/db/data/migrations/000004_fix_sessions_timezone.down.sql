-- Revert timezone changes for sessions table
-- Change TIMESTAMP WITH TIME ZONE columns back to TIMESTAMP

ALTER TABLE writers_sessions 
ALTER COLUMN expiry TYPE TIMESTAMP,
ALTER COLUMN created_at TYPE TIMESTAMP,
ALTER COLUMN updated_at TYPE TIMESTAMP,
ALTER COLUMN last_activity TYPE TIMESTAMP; 