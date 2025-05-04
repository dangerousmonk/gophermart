package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/service"
)

func (h *HTTPHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("LoginUser error on decoding body", slog.Any("err", err))
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.LoginUser(r.Context(), &req)
	if err != nil {
		slog.Error("Error on login user", slog.Any("error", err))

		switch {
		case errors.Is(err, service.ErrWrongPassword) || errors.Is(err, service.ErrNoUserFound):
			WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return

		default:
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	err = h.authenticator.SetAuth(user.ID, w, r)
	if err != nil {
		slog.Error("LoginUser error on setting cookies", slog.Any("error", err))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		slog.Error("LoginUser error on encoding response", slog.Any("error", err))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

}
