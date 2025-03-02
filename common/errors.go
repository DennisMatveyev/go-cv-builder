package common

import "errors"

var (
	ErrParseJSON = errors.New("Cannot parse JSON")
	ErrDatabase  = errors.New("Database error")
)
