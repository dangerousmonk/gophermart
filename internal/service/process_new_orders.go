package service

import (
	"context"
	"log/slog"
	"sync"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (s *GophermartService) ProccessPendingOrders(ctx context.Context, workerCount int) {
	orders, err := s.Repo.GetNewOrders(ctx)
	if err != nil {
		slog.Error("ProccessPendingOrders error on fetching orders from DB", slog.Any("error", err))
		return
	}
	if len(orders) == 0 {
		slog.Info("ProccessPendingOrders no orders found")
	}

	var wg sync.WaitGroup
	jobs := make(chan models.Order, len(orders))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for order := range jobs {
				resp, err := s.GetAccrual(order.Number)
				if err != nil {
					slog.Error("ProccessPendingOrders retrieving accrual from accrual system", slog.Any("error", err))
					continue
				}

				if resp == nil {
					slog.Info("ProccessPendingOrders skip order without response", slog.String("number", order.Number))
					continue
				}

				if order.Status != resp.Status {
					order.Status = resp.Status
				}

				if order.Accrual != resp.Accrual {
					order.Accrual = resp.Accrual
				}

				err = s.Repo.MakeAccrualToBalance(ctx, order)
				if err != nil {
					slog.Error("ProccessPendingOrders error on updating order", slog.Any("error", err))
					continue
				}
				slog.Info("ProccessPendingOrders update sucess", slog.String("number", order.Number))
			}
		}()
	}

	for _, order := range orders {
		jobs <- order
	}

	close(jobs)
	wg.Wait()
}
