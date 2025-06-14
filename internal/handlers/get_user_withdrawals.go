package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/service"
)

// GetWithdrawals godoc
//
//	@Summary		Get withdrawals
//	@Description	Get user withdrawals
//	@Accept			json
//
// @Param 		 Cookie header string  true "auth"     default(auth=xxx)
//
//	@Produce		json
//	@Tags			withdrawals
//	@Success		200 {object} models.Withdrawal
//	@Success		204
//	@Failure		401,500	{object}	errorResponse
//	@Router			/api/user/withdrawals   [get]
func (h *HTTPHandler) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		slog.Error("GetUserWithdrawals failed to cast userID", slog.Any("userID", userID))
		WriteErrorResponse(w, http.StatusUnauthorized, "No valid userID found")
		return
	}

	wds, err := h.service.GetUserWithdrawals(r.Context(), userID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoWithdrawals):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return

		default:
			slog.Error("GetUserWithdrawals error", slog.Any("error", err), slog.Int("userID", userID))
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
