package models

import "time"

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusRegistered OrderStatus = "REGISTERED"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int         `json:"id"`
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	UserID     int         `json:"user_id"`
	CreatedAt  time.Time   `json:"created_at"`
	ModifiedAt time.Time   `json:"modified_at"`
	Active     bool        `json:"active"`
	Accrual    float64     `json:"accrual"`
}
