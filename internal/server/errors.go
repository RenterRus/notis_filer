package server

import "errors"

var (
	ErrUserNotFound = errors.New("the user cannot be recognized")
)
