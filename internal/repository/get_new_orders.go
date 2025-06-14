package repository

import (
	"context"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) GetNewOrders(ctx context.Context) ([]models.Order, error) {
	const selectFields = "id,number,status,user_id,accrual,active,created_at,modified_at"
	const limit = 1000
	var orders []models.Order

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := r.conn.QueryContext(
		ctx,
		`SELECT `+selectFields+` FROM orders WHERE status IN ($1, $2, $3) LIMIT $4`,
		models.StatusNew, models.StatusRegistered, models.StatusProcessing, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ord models.Order
		if err = rows.Scan(&ord.ID, &ord.Number, &ord.Status, &ord.UserID, &ord.Accrual, &ord.Active, &ord.CreatedAt, &ord.ModifiedAt); err != nil {
			return nil, err
		}
		orders = append(orders, ord)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
