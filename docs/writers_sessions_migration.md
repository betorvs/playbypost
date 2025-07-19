# Writers Sessions Table Migration

## Purpose
This document outlines the migration process and structure of the `writers_sessions` table, which stores session information for writers in the Play-by-Post application. This table allows for persistent sessions, so writers are not logged out when the server restarts.

## Table Structure
The `writers_sessions` table has the following columns:

- `id`: bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY. A unique identifier for each session.
- `username`: VARCHAR(50) NOT NULL. The username of the writer.
- `token`: VARCHAR(255) NOT NULL UNIQUE. The session token.
- `user_id`: int NOT NULL REFERENCES writers(id) ON DELETE CASCADE. The ID of the writer.
- `expiry`: TIMESTAMP NOT NULL. The expiration date of the session.
- `client_type`: VARCHAR(20) NOT NULL. The type of client that initiated the session (e.g., 'browser', 'mobile', 'api', 'unknown').
- `client_info`: JSONB. Additional client context.
- `ip_address`: INET. The IP address of the client.
- `user_agent`: TEXT. The user agent of the client.
- `created_at`: TIMESTAMP NOT NULL DEFAULT NOW(). The timestamp when the session was created.
- `updated_at`: TIMESTAMP NOT NULL DEFAULT NOW(). The timestamp when the session was last updated.
- `last_activity`: TIMESTAMP NOT NULL DEFAULT NOW(). The timestamp of the last activity.

## Indexes
To improve performance, the `writers_sessions` table has the following indexes:

- `idx_writers_sessions_username`
- `idx_writers_sessions_token`
- `idx_writers_sessions_expiry`
- `idx_writers_sessions_user_id`
- `idx_writers_sessions_client_type`

## Foreign Key Constraints
To maintain data integrity, the `writers_sessions` table includes an `ON DELETE CASCADE` constraint for the `user_id` foreign key. This means that if a record in the `writers` table is deleted, all associated entries in the `writers_sessions` table for that `user_id` will be automatically deleted.

## Migration Files
The migration for the `writers_sessions` table is defined in the following files:

- `core/sys/db/data/migrations/000003_create_writers_sessions_table.up.sql`: Contains the SQL statements to create the `writers_sessions` table and its constraints.
- `core/sys/db/data/migrations/000003_create_writers_sessions_table.down.sql`: Contains the SQL statements to drop the `writers_sessions` table.

## Running the Migration
To apply the migration and create the `writers_sessions` table, use the `admin-ctl` command-line tool:

```bash
./admin-ctl db up
```

To roll back the migration and drop the `writers_sessions` table, use the `admin-ctl` command-line tool:

```bash
./admin-ctl db down
```
