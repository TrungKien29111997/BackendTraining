-- Drop trigger if exists
DROP TRIGGER IF EXISTS set_user_updated_at ON users;

-- Drop function if exists
DROP FUNCTION IF EXISTS update_user_updated_at_column;

-- Drop index if exists
DROP INDEX IF EXISTS user_email_status_idx;
DROP INDEX IF EXISTS user_level_idx;
DROP INDEX IF EXISTS user_status_idx;
DROP INDEX IF EXISTS user_created_at_idx;
DROP INDEX IF EXISTS user_deleted_at_idx;

-- Drop table
DROP TABLE IF EXISTS users;