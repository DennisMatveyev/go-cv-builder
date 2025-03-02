package auth

import "errors"

var (
	ErrUserExists         = errors.New("User with this email already exists")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrHashPassword       = errors.New("Failed to hash password")
	ErrSaveUser           = errors.New("Failed to save user")
	ErrGenerateToken      = errors.New("Could not generate token")
	ErrMissingAuthHeader  = errors.New("Missing Authorization header Bearer token")
	ErrInvalidToken       = errors.New("Invalid token")
)
