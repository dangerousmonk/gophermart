package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	internalErrors "github.com/dangerousmonk/gophermart/internal/errors"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/go-playground/validator/v10"
)

func (h *HTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("RegisterUser error on decoding body", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.service.RegisterUser(r.Context(), &req)
	if err != nil {
		slog.Error("Error on creating user", slog.Any("error", err))
		var validateErrs validator.ValidationErrors

		switch {
		case errors.Is(err, internalErrors.ErrLoginExists):
			http.Error(w, err.Error(), http.StatusConflict)
			return

		case errors.As(err, &validateErrs):
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = h.authenticator.SetAuth(userID, w, r)
	if err != nil {
		slog.Error("Error on setting cookies", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
