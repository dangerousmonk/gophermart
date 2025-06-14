package repository

import (
	"context"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) UploadOrder(ctx context.Context, orderNum string, userID int, status models.OrderStatus) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	res, err := r.conn.ExecContext(ctx, `INSERT INTO orders (number, user_id, status) VALUES ($1, $2, $3)`, orderNum, userID, status)
	if err != nil {
		return 0, err
	}
	rowsN, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsN, nil
}
