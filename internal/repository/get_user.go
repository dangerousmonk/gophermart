package repository

import (
	"context"
	"time"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (r *PostgresRepo) GetUser(ctx context.Context, login string) (models.User, error) {
	const selectFields = "id,login,password,active,created_at,modified_at,last_login_at"

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	row := r.conn.QueryRowContext(ctx, `SELECT `+selectFields+` FROM users WHERE login=$1`, login)

	var u models.User
	err := row.Scan(&u.ID, &u.Login, &u.Password, &u.Active, &u.CreatedAt, &u.ModifiedAt, &u.LastLoginAt)

	if err == nil {
		return u, nil
	}
	return u, err
}
