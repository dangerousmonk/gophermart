package service

import (
	"context"
	"log/slog"

	"github.com/dangerousmonk/gophermart/internal/models"
)

func (s *GophermartService) GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	wds, err := s.Repo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		slog.Error("GetUserWithdrawals failed to fetch orders", slog.Any("error", err))
		return nil, err
	}
	if len(wds) == 0 {
		return nil, ErrNoWithdrawals
	}
	return wds, nil
}
