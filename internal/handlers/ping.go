package handlers

import (
	"log/slog"
	"net/http"
)

func (h *HTTPHandler) Ping(w http.ResponseWriter, req *http.Request) {
	err := h.service.Ping(req.Context())
	if err != nil {
		slog.Error("Ping database unreachable ", slog.Any("err", err))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)

}
