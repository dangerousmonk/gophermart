package repository

import (
	"context"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) CreateUser(ctx context.Context, u *models.UserRequest) (int, error) {
	var userID int
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(
		ctx, `INSERT INTO users (login, password, last_login_at) VALUES ($1, $2, $3) RETURNING id`, u.Login, u.HashedPassword, time.Now(),
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO user_balance (user_id) VALUES ($1)`,
		userID,
	)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}
