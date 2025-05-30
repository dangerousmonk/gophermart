package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/middleware"
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
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		slog.Error("GetUserOrders failed to cast userID", slog.Any("userID", userID))
		WriteErrorResponse(w, http.StatusUnauthorized, "No valid userID found")
		return
	}
	orders, err := h.service.GetUserOrders(r.Context(), userID)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrNoOrders):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return

		default:
			slog.Error("GetUserOrders error", slog.Any("error", err), slog.Int("userID", userID))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		slog.Error("GetUserOrders error on encoding response", slog.Any("error", err), slog.Int("userID", userID))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

}
