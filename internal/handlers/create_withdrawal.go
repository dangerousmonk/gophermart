package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	internalErrors "github.com/dangerousmonk/gophermart/internal/errors"
	"github.com/dangerousmonk/gophermart/internal/models"
)

func (h *HTTPHandler) MakeWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req models.MakeWithdrawalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("MakeWithdrawal error on decoding body", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.service.MakeWithdrawal(r.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, internalErrors.ErrWrongOrderNum):
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		case errors.Is(err, internalErrors.ErrNoUserIDFound):
			http.Error(w, "User ID not found", http.StatusUnauthorized)
			return

		case errors.Is(err, internalErrors.ErrWithdrawalForOrderExists):
			http.Error(w, "Withdrawal for this order already registered", http.StatusConflict)
			return

		case errors.Is(err, internalErrors.ErrInsufficientBalance):
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return

		default:
			slog.Error("CreateWithdrawal error", slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}

}
