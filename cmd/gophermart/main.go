package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database.ApplyMigrations(cfg)
	db, err := database.InitDB(ctx, cfg.DatabaseDSN)
	if err != nil {
		logger.Error("InitDB failed", slog.Any("err", err))
	}
	defer db.Close()

	repo := repository.NewPostgresRepo(db)
	s := service.NewGophermartService(repo, cfg)

	app := server.NewGophermartApp(cfg, s)

	// server
	httpSrvDone := &sync.WaitGroup{}
	httpSrvDone.Add(1)
	srv := app.Start(httpSrvDone)

	// worker
	go s.StartAccrualWorker(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigChan
	logger.Info("Received os signal, shutting down")

	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeout)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
	}
}
