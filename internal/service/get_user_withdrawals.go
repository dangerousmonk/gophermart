package service

import (
	"context"
	"log/slog"

	appErrors "github.com/dangerousmonk/gophermart/internal/errors"
	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
)

func (s *GophermartService) GetUserWithdrawals(ctx context.Context) ([]models.Withdrawal, error) {
	id := ctx.Value(middleware.UserIDContextKey)
	if id == nil {
		slog.Error("GetUserWithdrawals no userID in context", slog.Any("error", id))
		return nil, appErrors.ErrNoUserIDFound
	}

	userID, ok := id.(int)
	if !ok {
		slog.Error("GetUserWithdrawals failed to cast userID", slog.Any("error", id))
		return nil, appErrors.ErrNoUserIDFound
	}

	wds, err := s.Repo.GetUserWithdrawals(ctx, userID)
	if err != nil {
		slog.Error("GetUserWithdrawals failed to fetch orders", slog.Any("error", id))
		return nil, err
	}
	if len(wds) == 0 {
		return nil, appErrors.ErrNoWithdrawals
	}
	return wds, nil
}
