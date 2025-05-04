package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dangerousmonk/gophermart/internal/service"
)

func (h *HTTPHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("UploadOrder error", slog.Any("error", err))
		http.Error(w, "Error on decoding body", http.StatusBadRequest)
		return
	}

	orderNum := strings.TrimSpace(string(body))
	_, err = h.service.UploadOrder(r.Context(), orderNum)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrWrongOrderNum):
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, service.ErrNoUserIDFound):
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return

		case errors.Is(err, service.ErrOrderExists):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return

		case errors.Is(err, service.ErrOrderExistsAnotherUser):
			http.Error(w, "Order uploaded by another user", http.StatusConflict)
			return

		default:
			slog.Error("UploadOrder error", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		return
	}

}
