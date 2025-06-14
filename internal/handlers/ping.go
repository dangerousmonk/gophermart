package handlers

import (
	"log/slog"
	"net/http"
)

// PingHandler godoc
// @Summary      DB Healthcheck
// @Description  Checks if DB is alive
// @Tags         Ping
// @Success      200
// @Failure      500
// @Router       /ping [get]
func (h *HTTPHandler) Ping(w http.ResponseWriter, req *http.Request) {
	err := h.service.Ping(req.Context())
	if err != nil {
		slog.Error("Ping database unreachable ", slog.Any("err", err))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

}
