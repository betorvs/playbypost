CREATE TABLE writer_users (
  id SERIAL PRIMARY KEY,
  writer_id INTEGER REFERENCES writers(id) ON DELETE CASCADE,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(writer_id, user_id)
);