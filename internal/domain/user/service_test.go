package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockRepository 모킹된 저장소 / Mocked repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) GetByID(id uint) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) GetByEmail(email string) (*User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) Update(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) List(query *ListUsersQuery) ([]*User, int64, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*User), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) Exists(id uint) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestService_Create(t *testing.T) {
	testCases := []struct {
		name          string
		request       *CreateUserRequest
		setupMock     func(*MockRepository)
		expectedError bool
		errorContains string
	}{
		{
			name: "successful user creation",
			request: &CreateUserRequest{
				Name:   "Test User",
				Email:  "test@example.com",
				Status: StatusActive,
			},
			setupMock: func(repo *MockRepository) {
				// Email doesn't exist
				repo.On("GetByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
				// Create succeeds
				repo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "email already exists",
			request: &CreateUserRequest{
				Name:   "Test User",
				Email:  "existing@example.com",
				Status: StatusActive,
			},
			setupMock: func(repo *MockRepository) {
				existingUser := &User{
					ID:     1,
					Name:   "Existing User",
					Email:  "existing@example.com",
					Status: StatusActive,
				}
				repo.On("GetByEmail", "existing@example.com").Return(existingUser, nil)
			},
			expectedError: true,
			errorContains: "email already exists",
		},
		{
			name: "database error during email check",
			request: &CreateUserRequest{
				Name:   "Test User",
				Email:  "test@example.com",
				Status: StatusActive,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("GetByEmail", "test@example.com").Return(nil, errors.New("database connection error"))
			},
			expectedError: true,
			errorContains: "failed to check email duplication",
		},
		{
			name: "database error during creation",
			request: &CreateUserRequest{
				Name:   "Test User",
				Email:  "test@example.com",
				Status: StatusActive,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("GetByEmail", "test@example.com").Return(nil, gorm.ErrRecordNotFound)
				repo.On("Create", mock.AnythingOfType("*user.User")).Return(errors.New("database insert error"))
			},
			expectedError: true,
			errorContains: "failed to create user",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tc.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			user, err := service.Create(tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.request.Name, user.Name)
				assert.Equal(t, tc.request.Email, user.Email)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	testCases := []struct {
		name          string
		userID        uint
		setupMock     func(*MockRepository)
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful user retrieval",
			userID: 1,
			setupMock: func(repo *MockRepository) {
				user := &User{
					ID:     1,
					Name:   "Test User",
					Email:  "test@example.com",
					Status: StatusActive,
				}
				repo.On("GetByID", uint(1)).Return(user, nil)
			},
			expectedError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			setupMock: func(repo *MockRepository) {
				repo.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
			errorContains: "user not found",
		},
		{
			name:   "database error",
			userID: 1,
			setupMock: func(repo *MockRepository) {
				repo.On("GetByID", uint(1)).Return(nil, errors.New("database connection error"))
			},
			expectedError: true,
			errorContains: "failed to get user",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tc.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			user, err := service.GetByID(tc.userID)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.userID, user.ID)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	newName := "Updated Name"
	newEmail := "updated@example.com"
	newStatus := StatusInactive

	testCases := []struct {
		name          string
		userID        uint
		request       *UpdateUserRequest
		setupMock     func(*MockRepository)
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful user update",
			userID: 1,
			request: &UpdateUserRequest{
				Name:   &newName,
				Email:  &newEmail,
				Status: &newStatus,
			},
			setupMock: func(repo *MockRepository) {
				existingUser := &User{
					ID:     1,
					Name:   "Old Name",
					Email:  "old@example.com",
					Status: StatusActive,
				}
				repo.On("GetByID", uint(1)).Return(existingUser, nil)
				repo.On("GetByEmail", "updated@example.com").Return(nil, gorm.ErrRecordNotFound)
				repo.On("Update", mock.AnythingOfType("*user.User")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			request: &UpdateUserRequest{
				Name: &newName,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
			},
			expectedError: true,
			errorContains: "user not found",
		},
		{
			name:   "email already exists",
			userID: 1,
			request: &UpdateUserRequest{
				Email: &newEmail,
			},
			setupMock: func(repo *MockRepository) {
				existingUser := &User{
					ID:     1,
					Name:   "Test User",
					Email:  "old@example.com",
					Status: StatusActive,
				}
				anotherUser := &User{
					ID:     2,
					Name:   "Another User",
					Email:  "updated@example.com",
					Status: StatusActive,
				}
				repo.On("GetByID", uint(1)).Return(existingUser, nil)
				repo.On("GetByEmail", "updated@example.com").Return(anotherUser, nil)
			},
			expectedError: true,
			errorContains: "email already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tc.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			user, err := service.Update(tc.userID, tc.request)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.userID, user.ID)
				if tc.request.Name != nil {
					assert.Equal(t, *tc.request.Name, user.Name)
				}
				if tc.request.Email != nil {
					assert.Equal(t, *tc.request.Email, user.Email)
				}
				if tc.request.Status != nil {
					assert.Equal(t, *tc.request.Status, user.Status)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	testCases := []struct {
		name          string
		userID        uint
		setupMock     func(*MockRepository)
		expectedError bool
		errorContains string
	}{
		{
			name:   "successful user deletion",
			userID: 1,
			setupMock: func(repo *MockRepository) {
				repo.On("Exists", uint(1)).Return(true, nil)
				repo.On("Delete", uint(1)).Return(nil)
			},
			expectedError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			setupMock: func(repo *MockRepository) {
				repo.On("Exists", uint(999)).Return(false, nil)
			},
			expectedError: true,
			errorContains: "user not found",
		},
		{
			name:   "database error during existence check",
			userID: 1,
			setupMock: func(repo *MockRepository) {
				repo.On("Exists", uint(1)).Return(false, errors.New("database connection error"))
			},
			expectedError: true,
			errorContains: "failed to check user existence",
		},
		{
			name:   "database error during deletion",
			userID: 1,
			setupMock: func(repo *MockRepository) {
				repo.On("Exists", uint(1)).Return(true, nil)
				repo.On("Delete", uint(1)).Return(errors.New("database delete error"))
			},
			expectedError: true,
			errorContains: "failed to delete user",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tc.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			err := service.Delete(tc.userID)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_List(t *testing.T) {
	testCases := []struct {
		name          string
		query         *ListUsersQuery
		setupMock     func(*MockRepository)
		expectedError bool
		errorContains string
		expectedCount int
		expectedTotal int64
	}{
		{
			name: "successful user listing",
			query: &ListUsersQuery{
				Offset: 0,
				Limit:  10,
			},
			setupMock: func(repo *MockRepository) {
				users := []*User{
					{ID: 1, Name: "User 1", Email: "user1@example.com", Status: StatusActive},
					{ID: 2, Name: "User 2", Email: "user2@example.com", Status: StatusActive},
				}
				repo.On("List", mock.AnythingOfType("*user.ListUsersQuery")).Return(users, int64(2), nil)
			},
			expectedError: false,
			expectedCount: 2,
			expectedTotal: 2,
		},
		{
			name: "database error during listing",
			query: &ListUsersQuery{
				Offset: 0,
				Limit:  10,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("List", mock.AnythingOfType("*user.ListUsersQuery")).Return(nil, int64(0), errors.New("database connection error"))
			},
			expectedError: true,
			errorContains: "failed to list users",
		},
		{
			name: "empty result set",
			query: &ListUsersQuery{
				Offset: 100,
				Limit:  10,
			},
			setupMock: func(repo *MockRepository) {
				repo.On("List", mock.AnythingOfType("*user.ListUsersQuery")).Return([]*User{}, int64(0), nil)
			},
			expectedError: false,
			expectedCount: 0,
			expectedTotal: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockRepository)
			tc.setupMock(mockRepo)
			service := NewService(mockRepo)

			// Execute
			users, total, err := service.List(tc.query)

			// Assert
			if tc.expectedError {
				assert.Error(t, err)
				assert.Nil(t, users)
				assert.Zero(t, total)
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, tc.expectedCount)
				assert.Equal(t, tc.expectedTotal, total)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

// 테스트 헬퍼 함수들 / Test helper functions

func createTestUser() *User {
	return &User{
		ID:     1,
		Name:   "Test User",
		Email:  "test@example.com",
		Status: StatusActive,
	}
}

func createTestCreateRequest() *CreateUserRequest {
	return &CreateUserRequest{
		Name:   "Test User",
		Email:  "test@example.com",
		Status: StatusActive,
	}
}

func createTestUpdateRequest() *UpdateUserRequest {
	name := "Updated Name"
	email := "updated@example.com"
	status := StatusInactive

	return &UpdateUserRequest{
		Name:   &name,
		Email:  &email,
		Status: &status,
	}
}

// 벤치마크 테스트 / Benchmark tests
func BenchmarkService_Create(b *testing.B) {
	mockRepo := new(MockRepository)
	mockRepo.On("GetByEmail", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)

	service := NewService(mockRepo)
	request := createTestCreateRequest()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		request.Email = "benchmark@example.com" // Unique email for each iteration
		service.Create(request)
	}
}

func BenchmarkService_GetByID(b *testing.B) {
	mockRepo := new(MockRepository)
	mockRepo.On("GetByID", mock.AnythingOfType("uint")).Return(createTestUser(), nil)

	service := NewService(mockRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetByID(1)
	}
}