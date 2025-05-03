package database

import (
	"errors"
	"log/slog"

	"github.com/dangerousmonk/gophermart/cmd/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ApplyMigrations(cfg *config.Config) {
	slog.Info("ApplyMigrations start", slog.String("dsn", cfg.DatabaseDSN))
	m, err := migrate.New("file://migrations/", cfg.DatabaseDSN)
	if err != nil {
		slog.Error("ApplyMigrations failed init instance ", slog.Any("err", err))
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("ApplyMigrations no change")
			return
		}
		slog.Error("ApplyMigrations failed", slog.Any("err", err))
	}
	slog.Info("ApplyMigrations success")
}
