package models

type AccrualExternal struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual"`
}
