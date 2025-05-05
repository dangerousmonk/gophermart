package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dangerousmonk/gophermart/internal/service"
)

func (h *HTTPHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("UploadOrder error on decoding body", slog.Any("error", err))
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	orderNum := strings.TrimSpace(string(body))
	if orderNum == "" {
		WriteErrorResponse(w, http.StatusBadRequest, "Missing order number")
		return
	}
	_, err = h.service.UploadOrder(r.Context(), orderNum)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrWrongOrderNum):
			WriteErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		case errors.Is(err, service.ErrNoUserIDFound):
			WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return

		case errors.Is(err, service.ErrOrderExists):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return

		case errors.Is(err, service.ErrOrderExistsAnotherUser):
			WriteErrorResponse(w, http.StatusConflict, err.Error())
			return

		default:
			slog.Error("UploadOrder error", slog.Any("error", err))
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		return
	}

}
