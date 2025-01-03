package repository

import "errors"

var (
	ErrsNotFound = errors.New("resource not found")
	ErrsConflict = errors.New("resource already exists")
)
