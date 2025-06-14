package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/go-playground/validator/v10"
)

// MakeWithdrawal godoc
//
//	@Summary		Make withdrawal
//	@Description	Make withdrawal
//	@Accept			json
//	@Produce		json
//	@Tags			withdrawals
//	@Param			data	body		models.MakeWithdrawalReq	true	"Request body"
//	@Success		200
//	@Failure		400,401,402,409,422,500	{object}	errorResponse
//	@Router			/api/user/balance/withdraw   [post]
func (h *HTTPHandler) MakeWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req models.MakeWithdrawalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("MakeWithdrawal error on decoding body", slog.Any("err", err))
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		slog.Error("MakeWithdrawal failed to cast userID", slog.Any("userID", userID))
		WriteErrorResponse(w, http.StatusUnauthorized, "No valid userID found")
		return
	}

	_, err := h.service.MakeWithdrawal(r.Context(), userID, req)

	if err != nil {
		var validateErrs validator.ValidationErrors

		switch {
		case errors.As(err, &validateErrs):
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return

		case errors.Is(err, service.ErrWrongOrderNum):
			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return

		case errors.Is(err, service.ErrWithdrawalForOrderExists):
			WriteErrorResponse(w, http.StatusConflict, err.Error())
			return

		case errors.Is(err, service.ErrInsufficientBalance):
			WriteErrorResponse(w, http.StatusPaymentRequired, err.Error())
			return

		default:
			slog.Error("CreateWithdrawal error", slog.Any("error", err), slog.Int("userID", userID))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}

}
