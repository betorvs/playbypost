# Writer Users Table Migration

## Purpose
This document outlines the migration process and structure of the `writer_users` table, which establishes a many-to-many relationship between writers and users in the Play-by-Post application. This table allows a single writer to be associated with multiple user accounts, enabling writers to manage characters across different player accounts.

## Table Structure
The `writer_users` table has the following columns:

- `id`: SERIAL PRIMARY KEY. A unique identifier for each writer-user association.
- `writer_id`: INTEGER, NOT NULL. This is a foreign key referencing the `id` column in the `writers` table. It identifies the writer associated with the user.
- `user_id`: INTEGER, NOT NULL. This is a foreign key referencing the `id` column in the `users` table. It identifies the user account associated with the writer.

The combination of `writer_id` and `user_id` forms a unique constraint, ensuring that a writer can only be associated with a specific user once.

## Foreign Key Constraints
To maintain data integrity, the `writer_users` table includes `ON DELETE CASCADE` constraints for both `writer_id` and `user_id` foreign keys. This means:

- If a record in the `writers` table is deleted, all associated entries in the `writer_users` table for that `writer_id` will be automatically deleted.
- If a record in the `users` table is deleted, all associated entries in the `writer_users` table for that `user_id` will be automatically deleted.

This cascading delete behavior prevents orphaned records and ensures consistency across the database.

## Migration Files
The migration for the `writer_users` table is defined in the following files:

- `core/sys/db/data/migrations/000002_create_writer_users_table.up.sql`: Contains the SQL statements to create the `writer_users` table and its constraints.
- `core/sys/db/data/migrations/000002_create_writer_users_table.down.sql`: Contains the SQL statements to drop the `writer_users` table.

## Running the Migration
To apply the migration and create the `writer_users` table, use the `admin-ctl` command-line tool:

```bash
go run app/admin-ctl/main.go db up
```

To roll back the migration and drop the `writer_users` table, use the `admin-ctl` command-line tool:

```bash
go run app/admin-ctl/main.go db down
```

## Verifying the Table Schema
To verify the schema of the `writer_users` table, including its columns and foreign key constraints, use the `admin-ctl` command-line tool:

```bash
go run app/admin-ctl/main.go db verify writer_users
```
