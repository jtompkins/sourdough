package auth

import (
	"time"
)

type User struct {
	Id        int       `json:"id" db:"id"`
	UserId    string    `json:"user_id" db:"user_id"`
	Provider  string    `json:"provider" db:"provider"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}