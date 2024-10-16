package repository

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource already exists")
	ErrDuplicateEmails    = errors.New("this email is already taken")
	ErrDuplicateUsernames = errors.New("this username is already taken")
)
