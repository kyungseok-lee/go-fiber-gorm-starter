package user

import (
	"time"

	"gorm.io/gorm"
)

// Status 사용자 상태 열거형 / User status enumeration
type Status string

const (
	StatusActive    Status = "active"
	StatusInactive  Status = "inactive"
	StatusSuspended Status = "suspended"
)

// User 사용자 모델 / User model
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Name      string         `json:"name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:255" validate:"required,email"`
	Status    Status         `json:"status" gorm:"not null;default:'active'" validate:"required,oneof=active inactive suspended"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 테이블 이름 지정 / Specify table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate 생성 전 훅 / Before create hook
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	// 기본 상태 설정 / Set default status
	if u.Status == "" {
		u.Status = StatusActive
	}
	return
}

// CreateUserRequest 사용자 생성 요청 구조체 / User creation request structure
type CreateUserRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=100"`
	Email  string `json:"email" validate:"required,email"`
	Status Status `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
}

// UpdateUserRequest 사용자 업데이트 요청 구조체 / User update request structure
type UpdateUserRequest struct {
	Name   *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Email  *string `json:"email,omitempty" validate:"omitempty,email"`
	Status *Status `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
}

// ListUsersQuery 사용자 목록 조회 쿼리 구조체 / User list query structure
type ListUsersQuery struct {
	Offset int    `query:"offset" validate:"min=0"`
	Limit  int    `query:"limit" validate:"min=1,max=100"`
	Status Status `query:"status" validate:"omitempty,oneof=active inactive suspended"`
	Search string `query:"search" validate:"omitempty,max=100"`
}

// Validate 쿼리 파라미터 검증 및 기본값 설정 / Validate query parameters and set defaults
func (q *ListUsersQuery) Validate() {
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
}

// ToUser CreateUserRequest를 User 모델로 변환 / Convert CreateUserRequest to User model
func (r *CreateUserRequest) ToUser() *User {
	user := &User{
		Name:  r.Name,
		Email: r.Email,
	}

	if r.Status != "" {
		user.Status = r.Status
	} else {
		user.Status = StatusActive
	}

	return user
}

// ApplyTo UpdateUserRequest를 기존 User 모델에 적용 / Apply UpdateUserRequest to existing User model
func (r *UpdateUserRequest) ApplyTo(user *User) {
	if r.Name != nil {
		user.Name = *r.Name
	}
	if r.Email != nil {
		user.Email = *r.Email
	}
	if r.Status != nil {
		user.Status = *r.Status
	}
}
