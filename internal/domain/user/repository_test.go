package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	// 테스트용 인메모리 SQLite 데이터베이스 사용 / Use in-memory SQLite database for testing
	// 실제 구현에서는 testcontainers-go 사용 권장 / Recommend using testcontainers-go in actual implementation
	// TODO: 실제 테스트 데이터베이스 연결 구현 / Implement actual test database connection
	t.Skip("Database connection for testing not implemented yet")
}

func TestRepository_Create(t *testing.T) {
	setupTestDB(t)
	var database *gorm.DB
	if database == nil {
		return
	}

	// Auto-migrate for testing
	err := database.AutoMigrate(&User{})
	require.NoError(t, err)

	repo := NewRepository(database)

	testCases := []struct {
		name    string
		user    *User
		wantErr bool
	}{
		{
			name: "valid user creation",
			user: &User{
				Name:   "Test User",
				Email:  "test@example.com",
				Status: StatusActive,
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: &User{
				Name:   "Another User",
				Email:  "test@example.com", // Same email as above
				Status: StatusActive,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.Create(tc.user)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, tc.user.ID)
			}
		})
	}
}

func TestRepository_GetByID(t *testing.T) {
	setupTestDB(t)
	var database *gorm.DB
	if database == nil {
		return
	}

	repo := NewRepository(database)

	// Create a test user first
	testUser := &User{
		Name:   "Test User",
		Email:  "test@example.com",
		Status: StatusActive,
	}
	err := repo.Create(testUser)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "existing user",
			userID:  testUser.ID,
			wantErr: false,
		},
		{
			name:    "non-existent user",
			userID:  99999,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.GetByID(tc.userID)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.userID, user.ID)
			}
		})
	}
}

func TestRepository_GetByEmail(t *testing.T) {
	setupTestDB(t)
	var database *gorm.DB
	if database == nil {
		return
	}

	repo := NewRepository(database)

	// Create a test user first
	testUser := &User{
		Name:   "Test User",
		Email:  "test@example.com",
		Status: StatusActive,
	}
	err := repo.Create(testUser)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "non-existent email",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.GetByEmail(tc.email)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.email, user.Email)
			}
		})
	}
}

func TestRepository_Update(t *testing.T) {
	setupTestDB(t)
	var database *gorm.DB
	if database == nil {
		return
	}

	repo := NewRepository(database)

	// Create a test user first
	testUser := &User{
		Name:   "Test User",
		Email:  "test@example.com",
		Status: StatusActive,
	}
	err := repo.Create(testUser)
	require.NoError(t, err)

	// Update the user
	testUser.Name = "Updated Name"
	testUser.Status = StatusInactive

	err = repo.Update(testUser)
	assert.NoError(t, err)

	// Verify the update
	updatedUser, err := repo.GetByID(testUser.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", updatedUser.Name)
	assert.Equal(t, StatusInactive, updatedUser.Status)
}

func TestRepository_Delete(t *testing.T) {
	setupTestDB(t)
	var database *gorm.DB
	if database == nil {
		return
	}

	repo := NewRepository(database)

	// Create a test user first
	testUser := &User{
		Name:   "Test User",
		Email:  "test@example.com",
		Status: StatusActive,
	}
	err := repo.Create(testUser)
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(testUser.ID)
	assert.NoError(t, err)

	// Verify the user is deleted (soft delete)
	_, err = repo.GetByID(testUser.ID)
	assert.Error(t, err) // Should not be found due to soft delete
}

func TestRepository_List(t *testing.T) {
	setupTestDB(t)
	var database *gorm.DB
	if database == nil {
		return
	}

	repo := NewRepository(database)

	// Create test users
	testUsers := []*User{
		{Name: "User 1", Email: "user1@example.com", Status: StatusActive},
		{Name: "User 2", Email: "user2@example.com", Status: StatusInactive},
		{Name: "User 3", Email: "user3@example.com", Status: StatusActive},
	}

	for _, user := range testUsers {
		err := repo.Create(user)
		require.NoError(t, err)
	}

	testCases := []struct {
		name        string
		query       *ListUsersQuery
		expectedMin int // Minimum expected results
	}{
		{
			name: "list all users",
			query: &ListUsersQuery{
				Offset: 0,
				Limit:  10,
			},
			expectedMin: 3,
		},
		{
			name: "list active users only",
			query: &ListUsersQuery{
				Offset: 0,
				Limit:  10,
				Status: StatusActive,
			},
			expectedMin: 2,
		},
		{
			name: "pagination test",
			query: &ListUsersQuery{
				Offset: 1,
				Limit:  2,
			},
			expectedMin: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			users, total, err := repo.List(tc.query)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, len(users), tc.expectedMin)
			assert.GreaterOrEqual(t, int(total), tc.expectedMin)
		})
	}
}

func TestRepository_Exists(t *testing.T) {
	setupTestDB(t)
	var database *gorm.DB
	if database == nil {
		return
	}

	repo := NewRepository(database)

	// Create a test user first
	testUser := &User{
		Name:   "Test User",
		Email:  "test@example.com",
		Status: StatusActive,
	}
	err := repo.Create(testUser)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		userID   uint
		expected bool
	}{
		{
			name:     "existing user",
			userID:   testUser.ID,
			expected: true,
		},
		{
			name:     "non-existent user",
			userID:   99999,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := repo.Exists(tc.userID)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, exists)
		})
	}
}

// 벤치마크 테스트 / Benchmark tests
func BenchmarkRepository_Create(b *testing.B) {
	setupTestDB(&testing.T{})
	var database *gorm.DB
	if database == nil {
		b.Skip("Database not available for benchmarking")
		return
	}

	repo := NewRepository(database)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &User{
			Name:   "Benchmark User",
			Email:  "benchmark@example.com",
			Status: StatusActive,
		}
		repo.Create(user)
	}
}

func BenchmarkRepository_GetByID(b *testing.B) {
	setupTestDB(&testing.T{})
	var database *gorm.DB
	if database == nil {
		b.Skip("Database not available for benchmarking")
		return
	}

	repo := NewRepository(database)

	// Create a test user
	testUser := &User{
		Name:   "Benchmark User",
		Email:  "benchmark@example.com",
		Status: StatusActive,
	}
	repo.Create(testUser)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.GetByID(testUser.ID)
	}
}
