-- Create users table
-- 사용자 테이블 생성

CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,  -- MySQL: AUTO_INCREMENT, PostgreSQL: SERIAL
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- PostgreSQL specific changes (comment/uncomment as needed):
-- For PostgreSQL, replace the above CREATE TABLE with:
-- CREATE TABLE users (
--     id BIGSERIAL PRIMARY KEY,
--     name VARCHAR(100) NOT NULL,
--     email VARCHAR(255) NOT NULL UNIQUE,
--     status VARCHAR(20) NOT NULL DEFAULT 'active',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     deleted_at TIMESTAMP WITH TIME ZONE NULL
-- );

-- Create indexes for better performance
-- 성능 향상을 위한 인덱스 생성
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Add check constraint for status
-- 상태 값 제약 조건 추가
ALTER TABLE users ADD CONSTRAINT chk_users_status 
CHECK (status IN ('active', 'inactive', 'suspended'));

-- PostgreSQL trigger for updated_at (uncomment for PostgreSQL):
-- CREATE OR REPLACE FUNCTION update_updated_at_column()
-- RETURNS TRIGGER AS $$
-- BEGIN
--     NEW.updated_at = CURRENT_TIMESTAMP;
--     RETURN NEW;
-- END;
-- $$ language 'plpgsql';
-- 
-- CREATE TRIGGER update_users_updated_at BEFORE UPDATE
--     ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();