package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) GetOrderByNumber(ctx context.Context, orderNum string) (models.Order, error) {
	const selectFields = "id,number,status,user_id,accrual,active,created_at,modified_at"

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	row := r.conn.QueryRowContext(ctx, `SELECT `+selectFields+` FROM orders WHERE number=$1`, orderNum)

	var ord models.Order
	err := row.Scan(&ord.ID, &ord.Number, &ord.Status, &ord.UserID, &ord.Accrual, &ord.Active, &ord.CreatedAt, &ord.ModifiedAt)

	if err == nil {
		return ord, nil
	}
	if err == sql.ErrNoRows {
		return ord, err
	}
	return ord, err
}

func (r *PostgresRepo) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	const selectFields = "id,number,status,user_id,accrual,active,created_at,modified_at"
	var orders []models.Order

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	rows, err := r.conn.QueryContext(ctx, `SELECT `+selectFields+` FROM orders WHERE user_id=$1 ORDER BY created_at DESC`, userID)
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
