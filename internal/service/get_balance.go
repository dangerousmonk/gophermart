package service

import (
	"context"
	"log/slog"

	appErrors "github.com/dangerousmonk/gophermart/internal/errors"
	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
)

func (s *GophermartService) GetBalance(ctx context.Context) (models.UserBalance, error) {
	var ub models.UserBalance
	id := ctx.Value(middleware.UserIDContextKey)
	if id == nil {
		slog.Error("GetBalance no userID in context", slog.Any("error", id))
		return ub, appErrors.ErrNoUserIDFound
	}

	userID, ok := id.(int)
	if !ok {
		slog.Error("GetBalance failed to cast userID", slog.Any("error", id))
		return ub, appErrors.ErrNoUserIDFound
	}

	ub, err := s.Repo.GetBalance(ctx, userID)
	if err != nil {
		slog.Error("GetBalance failed to fetch balance", slog.Any("error", id))
		return ub, err
	}
	return ub, nil
}
