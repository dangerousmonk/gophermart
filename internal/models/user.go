package models

import (
	"time"
)

type User struct {
	ID          int       `json:"id"`
	Login       string    `json:"login"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	LastLoginAt time.Time `json:"last_login_at"`
	Active      bool      `json:"active"`
}

type CreateUserReq struct {
	Login          string `json:"login" validate:"required,min=3,max=150"`
	Password       string `json:"password" validate:"required,min=5"`
	HashedPassword string `json:"-"`
}
