package repository

import (
	"context"
	"database/sql"

	"github.com/dangerousmonk/gophermart/internal/models"
)

type Repository interface {
	// Ping checks whether internal storage is up and running
	Ping(ctx context.Context) error
	// CreateUser creates new User
	CreateUser(ctx context.Context, u *models.CreateUserReq) (int, error)
	// GetUser searches user by login
	GetUser(ctx context.Context, login string) (models.User, error)
	// GetOrderByNumber searches order by order number
	GetOrderByNumber(ctx context.Context, orderNum string) (models.Order, error)
	// UploadOrder uploads new order
	UploadOrder(ctx context.Context, orderNum string, userID int, status models.OrderStatus) (int64, error)
	// GetUserOrders searches all orders by user
	GetUserOrders(ctx context.Context, userID int) ([]models.Order, error)
	// GetNewOrders searches all new orders registered in system
	GetNewOrders(ctx context.Context) ([]models.Order, error)
	// GetUserWithdrawals returns all withdrawals made by user
	GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error)
	// GetBalance returns user balance
	GetBalance(ctx context.Context, userID int) (models.UserBalance, error)
	// WithdrawFromBalance withdraws amount from user balance and makes entry in withdrawals table
	WithdrawFromBalance(ctx context.Context, orderNum string, userID int, amount float64) error
	// MakeAccrualToBalance updates order with accrual and updates user balance with accrual
	MakeAccrualToBalance(ctx context.Context, order models.Order) error
	// IsUniqueViolation check whether error is pg unique constraint error with specified constraint name
	IsUniqueViolation(err error, constraint string) bool
	// IsNoRows check whether error sql.ErrNoRows
	IsNoRows(err error) bool
}

type PostgresRepo struct {
	conn *sql.DB
}

func NewPostgresRepo(conn *sql.DB) *PostgresRepo {
	return &PostgresRepo{conn}
}
