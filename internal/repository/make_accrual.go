package repository

import (
	"context"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) MakeAccrualToBalance(ctx context.Context, order models.Order) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE orders SET status=$1, accrual=$2, modified_at=$3 WHERE number=$4`, order.Status, order.Accrual, time.Now(), order.Number)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE user_balance SET current=user_balance.current+$1 WHERE user_id=$2`,
		order.Accrual,
		order.UserID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
