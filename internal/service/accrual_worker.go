package service

import (
	"context"
	"log/slog"
	"time"
)

func (s *GophermartService) StartAccrualWorker(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.ProccessPendingOrders(ctx, 5)
		case <-ctx.Done():
			slog.Info("Accrual worker finished")
			return
		}

	}
}
