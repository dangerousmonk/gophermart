package handlers

import (
	"github.com/dangerousmonk/gophermart/internal/service"
	util "github.com/dangerousmonk/gophermart/internal/utils"
)

type HTTPHandler struct {
	service       service.GophermartService
	authenticator util.Authenticator
}

func NewHandler(s service.GophermartService, a util.Authenticator) *HTTPHandler {
	return &HTTPHandler{service: s, authenticator: a}
}
