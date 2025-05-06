package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/service"
)

// UserBalance godoc
//
//	@Summary		User balance
//	@Description	Get user balance
//	@Accept			json
//
// @Param 		 Cookie header string  true "auth"     default(auth=xxx)
//
//	@Produce		json
//	@Tags			balance
//	@Success		200 {object} models.UserBalance
//	@Failure		401,500	{object}	errorResponse
//	@Router			/api/user/balance   [get]
func (h *HTTPHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	ub, err := h.service.GetBalance(r.Context())

	if err != nil {
		switch {

		case errors.Is(err, service.ErrNoUserIDFound):
			WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return

		default:
			slog.Error("GetBalance error", slog.Any("error", err))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(ub); err != nil {
			slog.Error("GetUserWithdrawals error on encoding response", slog.Any("error", err))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

}
