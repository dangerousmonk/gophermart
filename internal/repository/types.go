package repository

import (
	"context"
	"database/sql"

	"github.com/dangerousmonk/gophermart/internal/models"
)

type Repository interface {
	// Ping checks whether internal storage is up and running
	Ping(ctx context.Context) error
	CreateUser(ctx context.Context, u *models.CreateUserReq) (int, error)
	GetUser(ctx context.Context, login string) (models.User, error)
	GetOrderByNumber(ctx context.Context, orderNum string) (models.Order, error)
	UploadOrder(ctx context.Context, orderNum string, userID int, status models.OrderStatus) (int64, error)
	GetUserOrders(ctx context.Context, userID int) ([]models.Order, error)
	GetNewOrders(ctx context.Context) ([]models.Order, error)
	GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error)
	GetBalance(ctx context.Context, userID int) (models.UserBalance, error)
	WithdrawFromBalance(ctx context.Context, orderNum string, userID int, amount float64) error
	MakeAccrualToBalance(ctx context.Context, order models.Order) error
	IsUniqueViolation(err error, constraint string) bool
	IsNoRows(err error) bool
}

type PostgresRepo struct {
	conn *sql.DB
}

func NewPostgresRepo(conn *sql.DB) *PostgresRepo {
	return &PostgresRepo{conn}
}
