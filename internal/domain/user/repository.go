package user

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Repository 사용자 저장소 인터페이스 / User repository interface
type Repository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(query *ListUsersQuery) ([]*User, int64, error)
	Exists(id uint) (bool, error)
}

// repository 사용자 저장소 구현체 / User repository implementation
type repository struct {
	db *gorm.DB
}

// NewRepository 새 사용자 저장소 생성 / Create new user repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create 사용자 생성 / Create user
func (r *repository) Create(user *User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID ID로 사용자 조회 / Get user by ID
func (r *repository) GetByID(id uint) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with id %d: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

// GetByEmail 이메일로 사용자 조회 / Get user by email
func (r *repository) GetByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with email %s: %w", email, err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// Update 사용자 업데이트 / Update user
func (r *repository) Update(user *User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// Delete 사용자 삭제 (소프트 삭제) / Delete user (soft delete)
func (r *repository) Delete(id uint) error {
	if err := r.db.Delete(&User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List 사용자 목록 조회 / List users
func (r *repository) List(query *ListUsersQuery) ([]*User, int64, error) {
	var users []*User
	var total int64

	// 기본 쿼리 / Base query
	db := r.db.Model(&User{})

	// 상태 필터링 / Status filtering
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 검색 필터링 (이름 또는 이메일) / Search filtering (name or email)
	if query.Search != "" {
		searchTerm := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ?", searchTerm, searchTerm)
	}

	// 총 개수 조회 / Get total count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 페이지네이션과 정렬 적용 / Apply pagination and sorting
	if err := db.Offset(query.Offset).
		Limit(query.Limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// Exists 사용자 존재 여부 확인 / Check if user exists
func (r *repository) Exists(id uint) (bool, error) {
	var count int64
	if err := r.db.Model(&User{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// WithTx 트랜잭션과 함께 저장소 반환 / Return repository with transaction
func (r *repository) WithTx(tx *gorm.DB) Repository {
	return &repository{db: tx}
}

// 향후 확장 가능한 메서드들 / Future extensible methods
// - BulkCreate: 대량 사용자 생성
// - BulkUpdate: 대량 사용자 업데이트  
// - GetActiveUsers: 활성 사용자만 조회
// - SearchByTags: 태그 기반 검색
// - GetUserStats: 사용자 통계 정보