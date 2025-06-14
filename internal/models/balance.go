package models

import "time"

type UserBalance struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Current   float64   `json:"current"`
	Withdrawn float64   `json:"withdrawn"`
	CreatedAt time.Time `json:"created_at"`
}
