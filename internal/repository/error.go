package repository

import "errors"

var (
	ErrsNotFound         = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("duplicate resource")
	ErrDuplicateUsername = errors.New("duplicate username")
)
