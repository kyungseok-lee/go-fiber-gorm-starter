package resp

import (
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse 에러 응답 구조체 / Error response structure
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail 에러 상세 정보 / Error detail information
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse 성공 응답 구조체 / Success response structure
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// PaginatedResponse 페이지네이션 응답 구조체 / Paginated response structure
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination 페이지네이션 정보 / Pagination information
type Pagination struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Total  int64 `json:"total"`
}

// Success 성공 응답 반환 / Return success response
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(SuccessResponse{Data: data})
}

// SuccessWithPagination 페이지네이션과 함께 성공 응답 반환 / Return success response with pagination
func SuccessWithPagination(c *fiber.Ctx, data interface{}, offset, limit int, total int64) error {
	return c.JSON(PaginatedResponse{
		Data: data,
		Pagination: Pagination{
			Offset: offset,
			Limit:  limit,
			Total:  total,
		},
	})
}

// Error 에러 응답 반환 / Return error response
func Error(c *fiber.Ctx, status int, code, message string, details ...interface{}) error {
	errResp := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}

	if len(details) > 0 {
		errResp.Error.Details = details[0]
	}

	return c.Status(status).JSON(errResp)
}

// BadRequest 400 에러 응답 / Return 400 error response
func BadRequest(c *fiber.Ctx, message string, details ...interface{}) error {
	return Error(c, fiber.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// Unauthorized 401 에러 응답 / Return 401 error response
func Unauthorized(c *fiber.Ctx, message string, details ...interface{}) error {
	return Error(c, fiber.StatusUnauthorized, "UNAUTHORIZED", message, details...)
}

// Forbidden 403 에러 응답 / Return 403 error response
func Forbidden(c *fiber.Ctx, message string, details ...interface{}) error {
	return Error(c, fiber.StatusForbidden, "FORBIDDEN", message, details...)
}

// NotFound 404 에러 응답 / Return 404 error response
func NotFound(c *fiber.Ctx, message string, details ...interface{}) error {
	return Error(c, fiber.StatusNotFound, "NOT_FOUND", message, details...)
}

// InternalServerError 500 에러 응답 / Return 500 error response
func InternalServerError(c *fiber.Ctx, message string, details ...interface{}) error {
	return Error(c, fiber.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, details...)
}

// Conflict 409 에러 응답 / Return 409 error response
func Conflict(c *fiber.Ctx, message string, details ...interface{}) error {
	return Error(c, fiber.StatusConflict, "CONFLICT", message, details...)
}

// UnprocessableEntity 422 에러 응답 / Return 422 error response
func UnprocessableEntity(c *fiber.Ctx, message string, details ...interface{}) error {
	return Error(c, fiber.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", message, details...)
}