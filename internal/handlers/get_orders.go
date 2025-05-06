package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/service"
)

// GetOrders godoc
//
//	@Summary		Get orders
//	@Description	Get user orders
//	@Accept			json
//
// @Param 		 Cookie header string  true "auth"     default(auth=xxx)
//
//	@Produce		json
//	@Tags			orders
//	@Success		200 {object} models.Order
//	@Failure		401,500	{object}	errorResponse
//	@Router			/api/user/orders   [get]
func (h *HTTPHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetUserOrders(r.Context())

	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoOrders):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return

		case errors.Is(err, service.ErrNoUserIDFound):
			slog.Error("GetUserOrders user ID not resolved", slog.Any("error", err))
			WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return

		default:
			slog.Error("GetUserOrders error", slog.Any("error", err))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			slog.Error("GetUserOrders error on encoding response", slog.Any("error", err))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		return
	}

}
