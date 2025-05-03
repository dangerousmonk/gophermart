package service

import (
	"context"

	internalErrors "github.com/dangerousmonk/gophermart/internal/errors"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/go-playground/validator/v10"
)

func (s *GophermartService) RegisterUser(ctx context.Context, req *models.CreateUserReq) (int, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(req)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return -1, errors
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return -1, err
	}
	req.HashedPassword = hashedPassword

	userID, err := s.Repo.CreateUser(ctx, req)
	if err != nil {
		if s.Repo.IsUniqueViolation(err, "users_login_key") {
			return -1, internalErrors.ErrLoginExists
		}
		return -1, err
	}
	return userID, nil
}
