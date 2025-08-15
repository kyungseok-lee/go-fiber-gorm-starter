-- Drop users table and related objects
-- 사용자 테이블 및 관련 객체 삭제

-- Drop PostgreSQL trigger if exists (uncomment for PostgreSQL):
-- DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
-- 인덱스 삭제
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_email;

-- Drop table
-- 테이블 삭제
DROP TABLE IF EXISTS users;