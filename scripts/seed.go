package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/config"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/db"
	"github.com/kyungseok-lee/go-fiber-gorm-starter/internal/domain/user"
)

func main() {
	// .env 파일 로드 / Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 설정 로드 / Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 개발환경에서만 실행 / Only run in development environment
	if !cfg.IsDev() {
		log.Fatalf("Seed script can only be run in development environment (ENV=dev or ENV=local)")
	}

	fmt.Println("🌱 Starting database seeding...")
	fmt.Printf("Environment: %s\n", cfg.Env)
	fmt.Printf("Database: %s\n", cfg.DBDriver)

	// 데이터베이스 연결 / Connect to database
	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}()

	// 기존 데이터 확인 / Check existing data
	var userCount int64
	if err := database.Model(&user.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Failed to count existing users: %v", err)
		return
	}

	if userCount > 0 {
		fmt.Printf("Users table already has %d records.\n", userCount)
		fmt.Print("Do you want to continue and add seed data? (y/N): ")

		var response string
		if _, err := fmt.Scanln(&response); err != nil {
			log.Printf("Failed to read input: %v", err)
			return
		}
		if response != "y" && response != "Y" {
			fmt.Println("Seed operation canceled.")
			os.Exit(0)
		}
	}

	// 시드 데이터 생성 / Create seed data
	seedUsers := []*user.User{
		{
			Name:   "John Doe",
			Email:  "john.doe@example.com",
			Status: user.StatusActive,
		},
		{
			Name:   "Jane Smith",
			Email:  "jane.smith@example.com",
			Status: user.StatusActive,
		},
		{
			Name:   "Bob Johnson",
			Email:  "bob.johnson@example.com",
			Status: user.StatusInactive,
		},
		{
			Name:   "Alice Brown",
			Email:  "alice.brown@example.com",
			Status: user.StatusActive,
		},
		{
			Name:   "Charlie Wilson",
			Email:  "charlie.wilson@example.com",
			Status: user.StatusSuspended,
		},
		{
			Name:   "Diana Davis",
			Email:  "diana.davis@example.com",
			Status: user.StatusActive,
		},
		{
			Name:   "Eve Miller",
			Email:  "eve.miller@example.com",
			Status: user.StatusActive,
		},
		{
			Name:   "Frank Garcia",
			Email:  "frank.garcia@example.com",
			Status: user.StatusInactive,
		},
		{
			Name:   "Grace Martinez",
			Email:  "grace.martinez@example.com",
			Status: user.StatusActive,
		},
		{
			Name:   "Henry Rodriguez",
			Email:  "henry.rodriguez@example.com",
			Status: user.StatusActive,
		},
	}

	// 트랜잭션으로 시드 데이터 삽입 / Insert seed data in transaction
	tx := database.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("Seed operation failed with panic: %v", r)
		}
	}()

	insertedCount := 0
	for _, seedUser := range seedUsers {
		// 이메일 중복 확인 / Check email duplication
		var existingUser user.User
		err := tx.Where("email = ?", seedUser.Email).First(&existingUser).Error
		if err == nil {
			fmt.Printf("⚠️  User with email %s already exists, skipping...\n", seedUser.Email)
			continue
		}

		// 사용자 생성 / Create user
		if err := tx.Create(seedUser).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Failed to create seed user %s: %v", seedUser.Email, err)
		}

		fmt.Printf("✅ Created user: %s (%s) - %s\n", seedUser.Name, seedUser.Email, seedUser.Status)
		insertedCount++
	}

	// 트랜잭션 커밋 / Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Failed to commit seed transaction: %v", err)
	}

	fmt.Printf("\n🎉 Seed operation completed!\n")
	fmt.Printf("Inserted %d new users out of %d total seed users.\n", insertedCount, len(seedUsers))

	// 최종 사용자 수 출력 / Print final user count
	var finalCount int64
	if err := database.Model(&user.User{}).Count(&finalCount).Error; err != nil {
		log.Printf("Warning: Failed to count final users: %v", err)
	} else {
		fmt.Printf("Total users in database: %d\n", finalCount)
	}
}
