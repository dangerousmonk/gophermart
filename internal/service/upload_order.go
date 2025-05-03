package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	appErrors "github.com/dangerousmonk/gophermart/internal/errors"
	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/utils"
)

func (s *GophermartService) UploadOrder(ctx context.Context, orderNum string) (models.Order, error) {
	var newOrder models.Order
	if !utils.IsValidOrderNumber(orderNum) {
		slog.Error("UploadOrder not valid order number", slog.Any("error", orderNum))
		return newOrder, appErrors.ErrWrongOrderNum
	}

	id := ctx.Value(middleware.UserIDContextKey)
	if id == nil {
		slog.Error("UploadOrder no userID in context", slog.Any("error", id))
		return newOrder, appErrors.ErrNoUserIDFound
	}

	userID, ok := id.(int)
	if !ok {
		slog.Error("UploadOrder failed to cast userID", slog.Any("error", id))
		return newOrder, appErrors.ErrNoUserIDFound
	}

	order, err := s.Repo.GetOrderByNumber(ctx, orderNum)

	switch {
	case err != nil && !errors.Is(err, sql.ErrNoRows):
		return newOrder, err

	case err != nil && errors.Is(err, sql.ErrNoRows):
		_, err := s.Repo.UploadOrder(ctx, orderNum, userID, models.StatusNew)
		if err != nil {
			slog.Error("UploadOrder failed to upload to postgres", slog.Any("error", err))
			return newOrder, err
		}
		return newOrder, nil

	case err == nil && order.UserID == userID:
		return newOrder, appErrors.ErrOrderExists

	case err == nil && order.UserID != userID:
		return newOrder, appErrors.ErrOrderExistsAnotherUser

	default:
		return newOrder, appErrors.ErrUnknown
	}
}
