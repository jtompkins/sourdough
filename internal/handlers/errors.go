package handlers

import "errors"

var ErrUnauthorized = errors.New("unauthorized")
var ErrUserNotFound = errors.New("user not found")
