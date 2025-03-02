package user

import "errors"

var (
	ErrInvalidParamFormat = errors.New("Invalid parameter format")
	ErrProfileExists      = errors.New("Profile already exists, update existing one")
)
