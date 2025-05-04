package service

import (
	"context"

	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/go-playground/validator/v10"
)

func (s *GophermartService) LoginUser(ctx context.Context, req *models.CreateUserReq) (models.User, error) {
	var user models.User
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(req)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return user, errors
	}

	user, err = s.Repo.GetUser(ctx, req.Login)
	if err != nil {
		if s.Repo.IsNoRows(err) {
			return user, ErrNoUserFound
		}
		return user, err
	}
	err = utils.CheckPassword(req.Password, user.Password)
	if err != nil {
		return user, ErrWrongPassword
	}
	return user, nil
}
