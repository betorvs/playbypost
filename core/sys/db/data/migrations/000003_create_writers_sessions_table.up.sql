CREATE TABLE writers_sessions (
  id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  username VARCHAR(50) NOT NULL,
  token VARCHAR(255) NOT NULL UNIQUE,
  user_id int NOT NULL REFERENCES writers(id) ON DELETE CASCADE,
  expiry TIMESTAMP NOT NULL,
  client_type VARCHAR(20) NOT NULL, -- 'browser', 'mobile', 'api', 'unknown'
  client_info JSONB, -- Additional client context
  ip_address INET,
  user_agent TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  last_activity TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_writers_sessions_username ON writers_sessions(username);
CREATE INDEX idx_writers_sessions_token ON writers_sessions(token);
CREATE INDEX idx_writers_sessions_expiry ON writers_sessions(expiry);
CREATE INDEX idx_writers_sessions_user_id ON writers_sessions(user_id);
CREATE INDEX idx_writers_sessions_client_type ON writers_sessions(client_type);