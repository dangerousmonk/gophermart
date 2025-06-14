package service

import (
	"context"
	"log/slog"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (s *GophermartService) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	orders, err := s.Repo.GetUserOrders(ctx, userID)
	if err != nil {
		slog.Error("GetUserOrders failed to fetch orders", slog.Any("error", err))
		return nil, err
	}
	if len(orders) == 0 {
		return nil, ErrNoOrders
	}
	return orders, nil
}
