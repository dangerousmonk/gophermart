package service

import (
	"context"
	"log/slog"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (s *GophermartService) GetBalance(ctx context.Context, userID int) (models.UserBalance, error) {
	ub, err := s.Repo.GetBalance(ctx, userID)
	if err != nil {
		slog.Error("GetBalance failed to fetch balance", slog.Any("error", err))
		return models.UserBalance{}, err
	}
	return ub, nil
}
