#!/bin/bash

# Database migration script
# 데이터베이스 마이그레이션 스크립트

set -e

# Load environment variables from .env file if it exists
# .env 파일이 있으면 환경변수 로드
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Default values
# 기본값 설정
DB_DRIVER=${DB_DRIVER:-mysql}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USER=${DB_USER:-user}
DB_PASS=${DB_PASS:-password}
DB_NAME=${DB_NAME:-fiber_gorm_starter}
DB_SSL_MODE=${DB_SSL_MODE:-disable}

# Build database URL based on driver
# 드라이버에 따른 데이터베이스 URL 구성
if [ "$DB_DRIVER" = "postgres" ]; then
    DATABASE_URL="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"
    DB_PORT=${DB_PORT:-5432}
elif [ "$DB_DRIVER" = "mysql" ]; then
    DATABASE_URL="mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"
else
    echo "❌ Unsupported database driver: $DB_DRIVER"
    echo "Supported drivers: mysql, postgres"
    exit 1
fi

MIGRATION_DIR="./migrations"

# Check if golang-migrate is installed
# golang-migrate 설치 확인
if ! command -v migrate &> /dev/null; then
    echo "❌ golang-migrate is not installed."
    echo "Please install it from: https://github.com/golang-migrate/migrate"
    echo ""
    echo "Installation options:"
    echo "1. Using go install:"
    echo "   go install -tags 'mysql postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    echo ""
    echo "2. Using homebrew (macOS):"
    echo "   brew install golang-migrate"
    echo ""
    echo "3. Using curl (Linux):"
    echo "   curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz"
    echo "   sudo mv migrate /usr/local/bin/"
    exit 1
fi

# Function to show usage
# 사용법 표시 함수
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  up [N]       Apply all or N up migrations"
    echo "  down [N]     Apply all or N down migrations"
    echo "  drop         Drop everything inside database"
    echo "  force N      Set version N but don't run migration (for fixing dirty state)"
    echo "  version      Print current migration version"
    echo "  create NAME  Create new migration files"
    echo "  status       Show migration status"
    echo ""
    echo "Environment variables:"
    echo "  DB_DRIVER    Database driver (mysql|postgres) [default: mysql]"
    echo "  DB_HOST      Database host [default: localhost]"
    echo "  DB_PORT      Database port [default: 3306 for mysql, 5432 for postgres]"
    echo "  DB_USER      Database user [default: user]"
    echo "  DB_PASS      Database password [default: password]"
    echo "  DB_NAME      Database name [default: fiber_gorm_starter]"
    echo "  DB_SSL_MODE  SSL mode for postgres [default: disable]"
    echo ""
    echo "Examples:"
    echo "  $0 up"
    echo "  $0 down 1"
    echo "  $0 create add_user_profile_table"
    echo "  DB_DRIVER=postgres $0 up"
}

# Function to check database connection
# 데이터베이스 연결 확인 함수
check_connection() {
    echo "🔍 Checking database connection..."
    echo "Driver: $DB_DRIVER"
    echo "Host: $DB_HOST:$DB_PORT"
    echo "Database: $DB_NAME"
    echo "User: $DB_USER"
    echo ""
    
    if [ "$DB_DRIVER" = "mysql" ]; then
        if command -v mysql &> /dev/null; then
            mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "SELECT 1;" "$DB_NAME" &> /dev/null
        else
            echo "⚠️  MySQL client not found, skipping connection test"
        fi
    elif [ "$DB_DRIVER" = "postgres" ]; then
        if command -v psql &> /dev/null; then
            PGPASSWORD="$DB_PASS" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" &> /dev/null
        else
            echo "⚠️  PostgreSQL client not found, skipping connection test"
        fi
    fi
}

# Parse command
# 명령어 파싱
COMMAND=${1:-}

case $COMMAND in
    "up")
        check_connection
        if [ -n "$2" ]; then
            echo "🚀 Applying $2 up migrations..."
            migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" up "$2"
        else
            echo "🚀 Applying all up migrations..."
            migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" up
        fi
        echo "✅ Migration completed!"
        ;;
    
    "down")
        check_connection
        if [ -n "$2" ]; then
            echo "⬇️  Applying $2 down migrations..."
            migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" down "$2"
        else
            echo "⬇️  Applying all down migrations..."
            migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" down
        fi
        echo "✅ Migration completed!"
        ;;
    
    "drop")
        check_connection
        echo "🗑️  Dropping all database objects..."
        echo "⚠️  This will destroy all data in the database!"
        read -p "Are you sure? (y/N): " confirm
        if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
            migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" drop
            echo "✅ Database dropped!"
        else
            echo "❌ Operation cancelled."
        fi
        ;;
    
    "force")
        if [ -z "$2" ]; then
            echo "❌ Version number required for force command"
            exit 1
        fi
        check_connection
        echo "🔧 Forcing version to $2..."
        migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" force "$2"
        echo "✅ Version forced to $2!"
        ;;
    
    "version")
        check_connection
        echo "📋 Current migration version:"
        migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" version
        ;;
    
    "create")
        if [ -z "$2" ]; then
            echo "❌ Migration name required"
            echo "Usage: $0 create MIGRATION_NAME"
            exit 1
        fi
        echo "📝 Creating new migration: $2"
        migrate create -ext sql -dir "$MIGRATION_DIR" "$2"
        echo "✅ Migration files created!"
        ;;
    
    "status")
        check_connection
        echo "📊 Migration status:"
        # Get current version
        VERSION=$(migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" version 2>/dev/null || echo "No version set")
        echo "Current version: $VERSION"
        
        # List migration files
        echo ""
        echo "Available migrations:"
        if [ -d "$MIGRATION_DIR" ]; then
            ls -la "$MIGRATION_DIR"/*.sql 2>/dev/null || echo "No migration files found"
        else
            echo "Migration directory not found: $MIGRATION_DIR"
        fi
        ;;
    
    "help"|"-h"|"--help"|"")
        show_usage
        ;;
    
    *)
        echo "❌ Unknown command: $COMMAND"
        echo ""
        show_usage
        exit 1
        ;;
esac