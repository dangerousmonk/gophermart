package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/service"
)

// UploadOrder godoc
//
//	@Summary		Upload order
//	@Description	Upload order
//	@Accept			plain
//	@Produce		json
//	@Tags			orders
//	@Param			order_number	body string	true	"Request body"
//	@Success		200
//	@Success		202
//	@Failure		400,401,409,422,500	{object}	errorResponse
//	@Router			/api/user/orders   [post]
func (h *HTTPHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("UploadOrder error on decoding body", slog.Any("error", err))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		slog.Error("UploadOrder failed to cast userID", slog.Any("userID", userID))
		WriteErrorResponse(w, http.StatusUnauthorized, "No valid userID found")
		return
	}

	orderNum := strings.TrimSpace(string(body))
	if orderNum == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing order number")
		return
	}
	_, err = h.service.UploadOrder(r.Context(), userID, orderNum)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrWrongOrderNum):
			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return

		case errors.Is(err, service.ErrOrderExists):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return

		case errors.Is(err, service.ErrOrderExistsAnotherUser):
			WriteErrorResponse(w, http.StatusConflict, err.Error())
			return

		default:
			slog.Error("UploadOrder error", slog.Any("error", err), slog.Int("userID", userID))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		return
	}

}
