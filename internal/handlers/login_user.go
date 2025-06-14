package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/go-playground/validator/v10"
)

// LoginUser godoc
//
//	@Summary		Login user
//	@Description	Login user
//	@Accept			json
//	@Produce		json
//	@Tags			users
//	@Param			data	body		models.UserRequest	true	"Request body"
//	@Success		200 {object} models.User
//	@Failure		400,401,500	{object}	errorResponse
//	@Router			/api/user/login   [post]
func (h *HTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("LoginUser error on decoding body", slog.Any("err", err))
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.LoginUser(r.Context(), &req)
	if err != nil {
		slog.Error("Error on login user", slog.Any("error", err))
		var validateErrs validator.ValidationErrors

		switch {
		case errors.Is(err, service.ErrWrongPassword) || errors.Is(err, service.ErrNoUserFound):
			WriteErrorResponse(w, http.StatusUnauthorized, utils.SanitizeError(err))
			return

		case errors.As(err, &validateErrs):
			WriteErrorResponse(w, http.StatusBadRequest, utils.SanitizeError(err))
			return

		default:
			WriteErrorResponse(w, http.StatusInternalServerError, utils.SanitizeError(err))
			return
		}
	}
	err = h.authenticator.SetAuth(user.ID, w, r)
	if err != nil {
		slog.Error("LoginUser error on setting cookies", slog.Any("error", utils.SanitizeError(err)), slog.Int("userID", user.ID))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.Error("LoginUser error on encoding response", slog.Any("error", utils.SanitizeError(err)), slog.Int("userID", user.ID))
		WriteErrorResponse(w, http.StatusInternalServerError, "Failed to login user")
		return
	}
	w.WriteHeader(http.StatusOK)

}
