package user

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Service 사용자 서비스 인터페이스 / User service interface
type Service interface {
	Create(req *CreateUserRequest) (*User, error)
	GetByID(id uint) (*User, error)
	Update(id uint, req *UpdateUserRequest) (*User, error)
	Delete(id uint) error
	List(query *ListUsersQuery) ([]*User, int64, error)
}

// service 사용자 서비스 구현체 / User service implementation
type service struct {
	repo Repository
}

// NewService 새 사용자 서비스 생성 / Create new user service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// Create 사용자 생성 / Create user
func (s *service) Create(req *CreateUserRequest) (*User, error) {
	logger := zap.L().With(zap.String("method", "user.service.Create"))

	// 이메일 중복 확인 / Check email duplication
	existingUser, err := s.repo.GetByEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to check email duplication", zap.Error(err), zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to check email duplication: %w", err)
	}
	
	if existingUser != nil {
		logger.Warn("Email already exists", zap.String("email", req.Email))
		return nil, fmt.Errorf("email already exists: %s", req.Email)
	}

	// 사용자 모델 생성 / Create user model
	user := req.ToUser()

	// 사용자 생성 / Create user
	if err := s.repo.Create(user); err != nil {
		logger.Error("Failed to create user", zap.Error(err), zap.String("email", req.Email))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("User created successfully", 
		zap.Uint("user_id", user.ID), 
		zap.String("email", user.Email))

	return user, nil
}

// GetByID ID로 사용자 조회 / Get user by ID
func (s *service) GetByID(id uint) (*User, error) {
	logger := zap.L().With(zap.String("method", "user.service.GetByID"), zap.Uint("user_id", id))

	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("User not found", zap.Uint("user_id", id))
			return nil, fmt.Errorf("user not found with id %d", id)
		}
		logger.Error("Failed to get user", zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update 사용자 업데이트 / Update user
func (s *service) Update(id uint, req *UpdateUserRequest) (*User, error) {
	logger := zap.L().With(
		zap.String("method", "user.service.Update"), 
		zap.Uint("user_id", id))

	// 기존 사용자 조회 / Get existing user
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("User not found for update", zap.Uint("user_id", id))
			return nil, fmt.Errorf("user not found with id %d", id)
		}
		logger.Error("Failed to get user for update", zap.Error(err))
		return nil, fmt.Errorf("failed to get user for update: %w", err)
	}

	// 이메일 중복 확인 (이메일이 변경되는 경우) / Check email duplication (if email is being changed)
	if req.Email != nil && *req.Email != user.Email {
		existingUser, err := s.repo.GetByEmail(*req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Failed to check email duplication for update", zap.Error(err))
			return nil, fmt.Errorf("failed to check email duplication: %w", err)
		}
		
		if existingUser != nil {
			logger.Warn("Email already exists for update", zap.String("email", *req.Email))
			return nil, fmt.Errorf("email already exists: %s", *req.Email)
		}
	}

	// 업데이트 요청 적용 / Apply update request
	req.ApplyTo(user)

	// 사용자 업데이트 / Update user
	if err := s.repo.Update(user); err != nil {
		logger.Error("Failed to update user", zap.Error(err))
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	logger.Info("User updated successfully", zap.Uint("user_id", user.ID))

	return user, nil
}

// Delete 사용자 삭제 / Delete user
func (s *service) Delete(id uint) error {
	logger := zap.L().With(
		zap.String("method", "user.service.Delete"), 
		zap.Uint("user_id", id))

	// 사용자 존재 확인 / Check user existence
	exists, err := s.repo.Exists(id)
	if err != nil {
		logger.Error("Failed to check user existence for delete", zap.Error(err))
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	
	if !exists {
		logger.Warn("User not found for delete", zap.Uint("user_id", id))
		return fmt.Errorf("user not found with id %d", id)
	}

	// 사용자 삭제 / Delete user
	if err := s.repo.Delete(id); err != nil {
		logger.Error("Failed to delete user", zap.Error(err))
		return fmt.Errorf("failed to delete user: %w", err)
	}

	logger.Info("User deleted successfully", zap.Uint("user_id", id))

	return nil
}

// List 사용자 목록 조회 / List users
func (s *service) List(query *ListUsersQuery) ([]*User, int64, error) {
	logger := zap.L().With(zap.String("method", "user.service.List"))

	// 쿼리 파라미터 검증 / Validate query parameters
	query.Validate()

	users, total, err := s.repo.List(query)
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	logger.Info("Users listed successfully", 
		zap.Int("count", len(users)), 
		zap.Int64("total", total),
		zap.Int("offset", query.Offset),
		zap.Int("limit", query.Limit))

	return users, total, nil
}

// 향후 확장 가능한 서비스 메서드들 / Future extensible service methods
// - CreateBatch: 대량 사용자 생성 (트랜잭션 내에서)
// - UpdateStatus: 사용자 상태 일괄 변경
// - SearchAdvanced: 고급 검색 기능
// - GetUserStatistics: 사용자 통계 정보
// - ActivateUser: 사용자 활성화
// - DeactivateUser: 사용자 비활성화