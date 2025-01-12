package repository

import "errors"

var (
	ErrsNotFound         = errors.New("resource not found")
	ErrsConflict         = errors.New("resource already exists")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)
