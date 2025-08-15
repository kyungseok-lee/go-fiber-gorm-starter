-- PostgreSQL initialization script
-- PostgreSQL 초기화 스크립트

-- Set timezone to Asia/Seoul
-- 타임존을 Asia/Seoul로 설정
SET timezone = 'Asia/Seoul';

-- Create database if not exists (already handled by docker-compose)
-- 데이터베이스 생성 (docker-compose에서 이미 처리됨)
-- This script runs after the database is created

-- Connect to the created database
-- 생성된 데이터베이스에 연결
\c fiber_gorm_starter;

-- Log initialization
-- 초기화 로그
SELECT 'PostgreSQL initialization completed' AS message;