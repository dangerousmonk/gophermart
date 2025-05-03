package handlers

import (
	"log/slog"
	"net/http"
)

func (h *HTTPHandler) Ping(w http.ResponseWriter, req *http.Request) {
	err := h.service.Ping(req.Context())
	if err != nil {
		slog.Error("Ping database unreachable ", slog.Any("err", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
