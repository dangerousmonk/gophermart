package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/service"
	util "github.com/dangerousmonk/gophermart/internal/utils"
)

type errorResponse struct {
	Error string `json:"error"`
}

type HTTPHandler struct {
	service       service.GophermartService
	authenticator util.Authenticator
}

func NewHandler(s service.GophermartService, a util.Authenticator) *HTTPHandler {
	return &HTTPHandler{service: s, authenticator: a}
}

func WriteErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}
