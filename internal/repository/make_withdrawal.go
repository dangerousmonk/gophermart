package repository

import (
	"context"
	"time"
)

func (r *PostgresRepo) WithdrawFromBalance(ctx context.Context, orderNum string, userID int, amount float64) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`UPDATE user_balance SET current=user_balance.current-$1, withdrawn=user_balance.withdrawn+$1 WHERE user_id=$2`,
		amount,
		userID,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO withdrawals (order_number, user_id, amount) VALUES ($1, $2, $3)`, orderNum, userID, amount)

	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
