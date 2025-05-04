package service

import (
	"context"
	"log/slog"

	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
)

func (s *GophermartService) GetUserOrders(ctx context.Context) ([]models.Order, error) {
	id := ctx.Value(middleware.UserIDContextKey)
	if id == nil {
		slog.Error("GetUserOrders no userID in context", slog.Any("error", id))
		return nil, ErrNoUserIDFound
	}

	userID, ok := id.(int)
	if !ok {
		slog.Error("GetUserOrders failed to cast userID", slog.Any("error", id))
		return nil, ErrNoUserIDFound
	}

	orders, err := s.Repo.GetUserOrders(ctx, userID)
	if err != nil {
		slog.Error("GetUserOrders failed to fetch orders", slog.Any("error", id))
		return nil, err
	}
	if len(orders) == 0 {
		return nil, ErrNoOrders
	}
	return orders, nil
}
