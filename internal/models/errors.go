package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid crednetials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
