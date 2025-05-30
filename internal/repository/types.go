package repository

import (
	"database/sql"
)

type PostgresRepo struct {
	conn *sql.DB
}

func NewPostgresRepo(conn *sql.DB) *PostgresRepo {
	return &PostgresRepo{conn}
}
