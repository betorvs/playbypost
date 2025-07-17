# Writer-User Association API Endpoints

This document outlines the backend logic and potential API endpoints for managing associations between writers and users in the Play-by-Post application. These functionalities are primarily handled within the `core/sys/db/pg` package.

## Database Interactions (`core/sys/db/pg/db_writer_users.go`)

The `db_writer_users.go` file contains the core functions for interacting with the `writer_users` table:

### `CreateWriterUserAssociation(ctx context.Context, writerID, userID int) (int, error)`
- **Purpose:** Creates a new association between a writer and a user.
- **Input:** `writerID` (int), `userID` (int).
- **Output:** The `id` (int) of the newly created association, or an error.
- **Validation:** Before insertion, this function performs the following checks:
    - Ensures `writerID` corresponds to an existing writer.
    - Ensures `userID` corresponds to an existing user.
    - Prevents duplicate associations between the same writer and user.
- **Potential API Endpoint:** `POST /api/writer_users` (Request body: `{ "writer_id": 123, "user_id": 456 }`)

### `GetUsersByWriterID(ctx context.Context, writerID int) ([]types.User, error)`
- **Purpose:** Retrieves all user records associated with a given writer.
- **Input:** `writerID` (int).
- **Output:** A slice of `types.User` structs, or an error.
- **Potential API Endpoint:** `GET /api/writers/{writerID}/users`

### `GetWritersByUserID(ctx context.Context, userID int) ([]types.Writer, error)`
- **Purpose:** Retrieves all writer records associated with a given user.
- **Input:** `userID` (int).
- **Output:** A slice of `types.Writer` structs, or an error.
- **Potential API Endpoint:** `GET /api/users/{userID}/writers`

### Validation Helper Functions

- `CheckWriterExists(ctx context.Context, writerID int) (bool, error)`
- `CheckUserExists(ctx context.Context, userID int) (bool, error)`
- `CheckWriterUserAssociationExists(ctx context.Context, writerID, userID int) (bool, error)`

These functions are used internally by `CreateWriterUserAssociation` to ensure data integrity before performing database operations.

## Error Handling

Errors returned by these functions will be standard Go `error` types or `fmt.Errorf` for custom messages, avoiding custom error types within the `pg` package.
