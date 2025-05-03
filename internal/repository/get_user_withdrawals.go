package repository

import (
	"context"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	const selectFields = "id,order_number,user_id,amount,created_at"
	var wds []models.Withdrawal

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	rows, err := r.conn.QueryContext(ctx, `SELECT `+selectFields+` FROM withdrawals WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var wd models.Withdrawal
		if err = rows.Scan(&wd.ID, &wd.OrderNumber, &wd.UserID, &wd.Amount, &wd.CreatedAt); err != nil {
			return nil, err
		}
		wds = append(wds, wd)

	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return wds, nil
}
