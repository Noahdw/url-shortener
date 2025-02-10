package service

import (
	"errors"
)

var (
	ErrInvalidURL = errors.New("invalid URL format")
	ErrNotFound   = errors.New("URL not found")
)
