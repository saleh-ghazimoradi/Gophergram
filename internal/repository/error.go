package repository

import "errors"

var (
	ErrsNotFound = errors.New("resource not found")
	ErrConflict  = errors.New("resource already exists")
)
