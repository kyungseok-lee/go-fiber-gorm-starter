// Package user provides user domain logic, handlers, and data models
package user

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/kyungseok-lee/go-fiber-gorm-starter/pkg/resp"
)

const (
	errEmailAlreadyExists = "email already exists"
	errUserNotFound       = "user not found"
)

// Handler 사용자 HTTP 핸들러 / User HTTP handler
type Handler struct {
	service Service
}

// NewHandler 새 사용자 핸들러 생성 / Create new user handler
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Create 사용자 생성 / Create user
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User creation request"
// @Success 201 {object} resp.SuccessResponse{data=User}
// @Failure 400 {object} resp.ErrorResponse
// @Failure 409 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /v1/users [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateUserRequest

	// 요청 바디 파싱 / Parse request body
	if err := c.BodyParser(&req); err != nil {
		return resp.BadRequest(c, "Invalid request body", err.Error())
	}

	// 기본 필드 검증 / Basic field validation
	if req.Name == "" {
		return resp.BadRequest(c, "Name is required")
	}
	if req.Email == "" {
		return resp.BadRequest(c, "Email is required")
	}
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return resp.BadRequest(c, "Name must be between 2 and 100 characters")
	}

	// TODO: 더 정교한 검증 로직 추가 가능 / Can add more sophisticated validation logic
	// - Email 형식 검증 (정규표현식)
	// - 비밀번호 강도 검증 (향후 추가 시)
	// - 사용자 정의 검증 규칙

	user, err := h.service.Create(&req)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			(err.Error() != "" && (err.Error() == errEmailAlreadyExists ||
				(len(err.Error()) > 20 && err.Error()[:20] == errEmailAlreadyExists))) {
			return resp.Conflict(c, "Email already exists")
		}
		zap.L().Error("Failed to create user", zap.Error(err))
		return resp.InternalServerError(c, "Failed to create user")
	}

	return c.Status(fiber.StatusCreated).JSON(resp.SuccessResponse{Data: user})
}

// GetByID ID로 사용자 조회 / Get user by ID
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} resp.SuccessResponse{data=User}
// @Failure 400 {object} resp.ErrorResponse
// @Failure 404 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /v1/users/{id} [get]
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return resp.BadRequest(c, "Invalid user ID")
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			(err.Error() != "" && (err.Error() == errUserNotFound ||
				len(err.Error()) > 15 && err.Error()[:15] == errUserNotFound)) {
			return resp.NotFound(c, "User not found")
		}
		zap.L().Error("Failed to get user", zap.Error(err), zap.Uint64("user_id", id))
		return resp.InternalServerError(c, "Failed to get user")
	}

	return resp.Success(c, user)
}

// Update 사용자 업데이트 / Update user
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User update request"
// @Success 200 {object} resp.SuccessResponse{data=User}
// @Failure 400 {object} resp.ErrorResponse
// @Failure 404 {object} resp.ErrorResponse
// @Failure 409 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /v1/users/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return resp.BadRequest(c, "Invalid user ID")
	}

	var req UpdateUserRequest
	if parseErr := c.BodyParser(&req); parseErr != nil {
		return resp.BadRequest(c, "Invalid request body", parseErr.Error())
	}

	// 기본 필드 검증 / Basic field validation
	if req.Name != nil && (*req.Name == "" || len(*req.Name) < 2 || len(*req.Name) > 100) {
		return resp.BadRequest(c, "Name must be between 2 and 100 characters")
	}
	if req.Email != nil && *req.Email == "" {
		return resp.BadRequest(c, "Email cannot be empty")
	}

	user, err := h.service.Update(uint(id), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			(err.Error() != "" && (err.Error() == errUserNotFound ||
				len(err.Error()) > 15 && err.Error()[:15] == errUserNotFound)) {
			return resp.NotFound(c, "User not found")
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			(err.Error() != "" && (err.Error() == errEmailAlreadyExists ||
				(len(err.Error()) > 20 && err.Error()[:20] == errEmailAlreadyExists))) {
			return resp.Conflict(c, "Email already exists")
		}
		zap.L().Error("Failed to update user", zap.Error(err), zap.Uint64("user_id", id))
		return resp.InternalServerError(c, "Failed to update user")
	}

	return resp.Success(c, user)
}

// Delete 사용자 삭제 / Delete user
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} resp.ErrorResponse
// @Failure 404 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /v1/users/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return resp.BadRequest(c, "Invalid user ID")
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) ||
			(err.Error() != "" && (err.Error() == errUserNotFound ||
				len(err.Error()) > 15 && err.Error()[:15] == errUserNotFound)) {
			return resp.NotFound(c, "User not found")
		}
		zap.L().Error("Failed to delete user", zap.Error(err), zap.Uint64("user_id", id))
		return resp.InternalServerError(c, "Failed to delete user")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// List 사용자 목록 조회 / List users
// @Summary List users
// @Description Get list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(20)
// @Param status query string false "Filter by status" Enums(active, inactive, suspended)
// @Param search query string false "Search by name or email"
// @Success 200 {object} resp.PaginatedResponse{data=[]User}
// @Failure 400 {object} resp.ErrorResponse
// @Failure 500 {object} resp.ErrorResponse
// @Router /v1/users [get]
func (h *Handler) List(c *fiber.Ctx) error {
	var query ListUsersQuery

	// 쿼리 파라미터 파싱 / Parse query parameters
	if err := c.QueryParser(&query); err != nil {
		return resp.BadRequest(c, "Invalid query parameters", err.Error())
	}

	// 쿼리 검증 및 기본값 설정 / Validate query and set defaults
	query.Validate()

	users, total, err := h.service.List(&query)
	if err != nil {
		zap.L().Error("Failed to list users", zap.Error(err))
		return resp.InternalServerError(c, "Failed to list users")
	}

	return resp.SuccessWithPagination(c, users, query.Offset, query.Limit, total)
}

// 향후 확장 가능한 핸들러 메서드들 / Future extensible handler methods
// - BulkCreate: 대량 사용자 생성
// - BulkUpdate: 대량 사용자 업데이트
// - Export: 사용자 데이터 내보내기 (CSV, Excel 등)
// - Import: 사용자 데이터 가져오기
// - GetProfile: 사용자 프로필 조회 (확장된 정보)
// - UpdateStatus: 사용자 상태만 변경
