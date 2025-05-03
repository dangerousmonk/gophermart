package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	internalErrors "github.com/dangerousmonk/gophermart/internal/errors"
)

func (h *HTTPHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetUserOrders(r.Context())

	if err != nil {
		switch {
		case errors.Is(err, internalErrors.ErrNoOrders):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return

		case errors.Is(err, internalErrors.ErrNoUserIDFound):
			slog.Error("GetUserOrders user ID not resolved", slog.Any("error", err))
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return

		default:
			slog.Error("GetUserOrders error", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			slog.Error("GetUserOrders error on encoding response", slog.Any("error", err))
			http.Error(w, `{"error":" failed, to encode response"}`, http.StatusInternalServerError)
			return
		}
		return
	}

}
