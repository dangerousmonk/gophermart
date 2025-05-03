package repository

import (
	"database/sql"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *PostgresRepo) IsUniqueViolation(err error, constraint string) bool {
	if pgError, ok := err.(*pgconn.PgError); ok {
		return pgError.Code == pgerrcode.UniqueViolation && pgError.ConstraintName == constraint
	}
	return false
}

func (r *PostgresRepo) IsNoRows(err error) bool {
	return err == sql.ErrNoRows
}
