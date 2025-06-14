package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dangerousmonk/gophermart/internal/models"
	"github.com/dangerousmonk/gophermart/internal/service"
	util "github.com/dangerousmonk/gophermart/internal/utils"
)

//go:generate mockgen -package mocks -source types.go -destination ./mocks/mock_service.go Service
type Service interface {
	// StartAccrualWorker starts background worker to process new orders
	StartAccrualWorker(ctx context.Context)
	// GetAccrual fetches accrual info by orderNumber using accrual application
	GetAccrual(orderNumber string) (*models.AccrualExternal, error)
	// GetAccrual fetches current user balance
	GetBalance(ctx context.Context, userID int) (models.UserBalance, error)
	// GetUserOrders fetches orders uploaded by user
	GetUserOrders(ctx context.Context, userID int) ([]models.Order, error)
	// UploadOrder persists new order by given user
	UploadOrder(ctx context.Context, userID int, orderNum string) (models.Order, error)
	// ProccessPendingOrders fetches new orders and sync their information with accrual application
	ProccessPendingOrders(ctx context.Context)
	// GetUserWithdrawals fetches withdrawals made by user
	GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error)
	// MakeWithdrawal withdraws funds from user current balance and saves history
	MakeWithdrawal(ctx context.Context, userID int, wdReq models.MakeWithdrawalReq) (models.Withdrawal, error)
	// LoginUser checks user credentials and authenticate him
	LoginUser(ctx context.Context, req *models.UserRequest) (models.User, error)
	// RegisterUser persists new user to database
	RegisterUser(ctx context.Context, req *models.UserRequest) (int, error)
	// Ping checks whether database is ok
	Ping(ctx context.Context) error
}

type HTTPHandler struct {
	service       Service
	authenticator util.Authenticator
}

func NewHandler(s service.GophermartService, a util.Authenticator) *HTTPHandler {
	return &HTTPHandler{service: &s, authenticator: a}
}

type errorResponse struct {
	Error string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}
