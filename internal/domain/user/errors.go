package user

import "errors"

var (
	// ErrEmailAlreadyExists is returned when a user email is already registered.
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrUserNotFound is returned when a user lookup cannot find a matching row.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidStatus is returned when a user status is outside the supported enum values.
	ErrInvalidStatus = errors.New("invalid user status")
)
