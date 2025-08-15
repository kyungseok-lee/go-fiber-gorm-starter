-- MySQL initialization script
-- MySQL 초기화 스크립트

-- Set timezone to Asia/Seoul
-- 타임존을 Asia/Seoul로 설정
SET GLOBAL time_zone = '+09:00';
SET SESSION time_zone = '+09:00';

-- Create database if not exists (already handled by docker-compose)
-- 데이터베이스 생성 (docker-compose에서 이미 처리됨)
-- CREATE DATABASE IF NOT EXISTS fiber_gorm_starter CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
-- 데이터베이스 사용
USE fiber_gorm_starter;

-- Grant privileges to the user (already handled by docker-compose)
-- 사용자에게 권한 부여 (docker-compose에서 이미 처리됨)
-- GRANT ALL PRIVILEGES ON fiber_gorm_starter.* TO 'user'@'%';
-- FLUSH PRIVILEGES;

-- Log initialization
-- 초기화 로그
SELECT 'MySQL initialization completed' AS message;