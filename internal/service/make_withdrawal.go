package service

import (
	"context"
	"log/slog"

	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/go-playground/validator/v10"
)

func (s *GophermartService) MakeWithdrawal(ctx context.Context, wdReq models.MakeWithdrawalReq) (models.Withdrawal, error) {
	var wd models.Withdrawal

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(wdReq)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return wd, errors
	}

	if !utils.IsValidOrderNumber(wdReq.Order) {
		slog.Error("CreateWithdrawal not valid order number", slog.Any("error", wdReq.Order))
		return wd, ErrWrongOrderNum
	}

	id := ctx.Value(middleware.UserIDContextKey)
	if id == nil {
		slog.Error("CreateWithdrawal no userID in context", slog.Any("error", id))
		return wd, ErrNoUserIDFound
	}

	userID, ok := id.(int)
	if !ok {
		slog.Error("CreateWithdrawal failed to cast userID", slog.Any("error", id))
		return wd, ErrNoUserIDFound
	}

	balance, err := s.Repo.GetBalance(ctx, userID)
	if err != nil {
		slog.Error("CreateWithdrawal failed to check balance", slog.Any("error", id))
		return wd, err
	}
	if balance.Current < wdReq.Sum {
		slog.Error("CreateWithdrawal insufficient funds", slog.Any("error", id))
		return wd, ErrInsufficientBalance
	}

	err = s.Repo.WithdrawFromBalance(ctx, wdReq.Order, userID, wdReq.Sum)
	if err != nil {
		if s.Repo.IsUniqueViolation(err, "withdrawals_order_number_key") {
			return wd, ErrWithdrawalForOrderExists
		}
		return wd, err
	}
	return wd, nil
}
