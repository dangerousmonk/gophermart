package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	internalErrors "github.com/dangerousmonk/gophermart/internal/errors"
	"github.com/dangerousmonk/gophermart/internal/models"
)

func (h *HTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("LoginUser error on decoding body", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.LoginUser(r.Context(), &req)
	if err != nil {
		slog.Error("Error on login user", slog.Any("error", err))

		switch {
		case errors.Is(err, internalErrors.ErrWrongPassword) || errors.Is(err, internalErrors.ErrNoUserFound):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = h.authenticator.SetAuth(user.ID, w, r)
	if err != nil {
		slog.Error("Error on setting cookies", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error on encoding response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
