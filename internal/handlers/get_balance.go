package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/service"
)

func (h *HTTPHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ub, err := h.service.GetBalance(r.Context())

	if err != nil {
		switch {

		case errors.Is(err, service.ErrNoUserIDFound):
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return

		default:
			slog.Error("GetBalance error", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(ub); err != nil {
			slog.Error("GetUserWithdrawals error on encoding response", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

}
