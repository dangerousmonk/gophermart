package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) GetBalance(ctx context.Context, userID int) (models.UserBalance, error) {
	const selectFields = "id,user_id,current,withdrawn,created_at"

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	row := r.conn.QueryRowContext(ctx, `SELECT `+selectFields+` FROM user_balance WHERE user_id=$1`, userID)

	var ub models.UserBalance
	err := row.Scan(&ub.ID, &ub.UserID, &ub.Current, &ub.Withdrawn, &ub.CreatedAt)

	if err == nil {
		return ub, nil
	}
	if err == sql.ErrNoRows {
		return ub, err
	}
	return ub, err
}
