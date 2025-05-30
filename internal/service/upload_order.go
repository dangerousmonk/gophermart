package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/utils"
)

func (s *GophermartService) UploadOrder(ctx context.Context, userID int, orderNum string) (models.Order, error) {
	var newOrder models.Order
	if !utils.IsValidOrderNumber(orderNum) {
		slog.Error("UploadOrder not valid order number", slog.Any("error", orderNum))
		return newOrder, ErrWrongOrderNum
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
		return newOrder, ErrOrderExists

	case err == nil && order.UserID != userID:
		return newOrder, ErrOrderExistsAnotherUser

	default:
		return newOrder, ErrUnknown
	}
}
