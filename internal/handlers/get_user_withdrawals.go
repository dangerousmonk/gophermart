package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/service"
)

func (h *HTTPHandler) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	wds, err := h.service.GetUserWithdrawals(r.Context())

	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoWithdrawals):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return

		case errors.Is(err, service.ErrNoUserIDFound):
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return

		default:
			slog.Error("GetUserWithdrawals error", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(wds); err != nil {
			slog.Error("GetUserWithdrawals error on encoding response", slog.Any("error", err))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

}
