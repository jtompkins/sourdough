package shared

import "errors"

var ErrUnauthorized = errors.New("unauthorized")
var ErrUserNotFound = errors.New("user not found")

type UserInfo struct {
	Id       int
	UserId   string
	Provider string
}