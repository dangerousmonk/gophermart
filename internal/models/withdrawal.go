package models

import "time"

type Withdrawal struct {
	ID          int       `json:"id"`
	OrderNumber string    `json:"order"`
	Amount      float64   `json:"sum"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type MakeWithdrawalReq struct {
	Order string  `json:"order" validate:"required,min=3"`
	Sum   float64 `json:"sum" validate:"required,gt=0"`
}
