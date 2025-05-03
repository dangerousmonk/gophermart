package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/dangerousmonk/gophermart/cmd/config"
	"github.com/dangerousmonk/gophermart/internal/database"
	"github.com/dangerousmonk/gophermart/internal/repository"
	"github.com/dangerousmonk/gophermart/internal/server"
	"github.com/dangerousmonk/gophermart/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg := config.InitConfig()
	ctx := context.Background()

	database.ApplyMigrations(cfg)
	db, err := database.InitDB(ctx, cfg.DatabaseDSN)
	if err != nil {
		logger.Error("InitDB failed", slog.Any("err", err))
	}
	defer db.Close()

	repo := repository.NewPostgresRepo(db)
	s := service.NewGophermartService(repo, cfg)

	server := server.NewGophermartApp(cfg, s)

	go s.StartAccrualWorker(ctx)
	err = server.Start()

	if err != nil {
		logger.Error("Server start failed", slog.Any("err", err))
	}
}
