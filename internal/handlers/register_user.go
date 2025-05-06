package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/go-playground/validator/v10"
)

// RegisterUser godoc
//
//	@Summary		Register user
//	@Description	Register user
//	@Accept			json
//	@Produce		json
//	@Tags			users
//	@Param			data	body		models.UserRequest	true	"Request body"
//	@Success		200
//	@Failure		400,409,500	{object}	errorResponse
//	@Router			/api/user/register   [post]
func (h *HTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("RegisterUser error on decoding body", slog.Any("error", err))
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := h.service.RegisterUser(r.Context(), &req)
	if err != nil {
		slog.Error("RegisterUser error on creating user", slog.Any("error", err))
		var validateErrs validator.ValidationErrors

		switch {
		case errors.Is(err, service.ErrLoginExists):
			WriteErrorResponse(w, http.StatusConflict, err.Error())
			return

		case errors.As(err, &validateErrs):
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return

		default:
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	err = h.authenticator.SetAuth(userID, w, r)
	if err != nil {
		slog.Error("RegisterUser error on setting cookies", slog.Any("error", err))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
