package server

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/dangerousmonk/gophermart/cmd/config"
	"github.com/dangerousmonk/gophermart/internal/handlers"
	appMiddleware "github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/go-chi/chi/v5"
)

type GophermartApp struct {
	Config  *config.Config
	Service *service.GophermartService
}

func NewGophermartApp(cfg *config.Config, s *service.GophermartService) *GophermartApp {
	return &GophermartApp{
		Config:  cfg,
		Service: s,
	}
}

func (app *GophermartApp) Start(wg *sync.WaitGroup) *http.Server {
	r := app.initRouter()
	srv := &http.Server{Addr: app.Config.ServerAddr}
	go func() {
		defer wg.Done()

		if err := http.ListenAndServe(app.Config.ServerAddr, r); err != http.ErrServerClosed {
			slog.Error("Server failed initialize", slog.String("address", app.Config.ServerAddr))
		}
	}()
	return srv
}

func (app *GophermartApp) initRouter() *chi.Mux {
	r := chi.NewRouter()
	jwtAuthenticator, err := utils.NewJWTAuthenticator(app.Config.JWTSecret)
	if err != nil {
		slog.Error("Server failed initialize jwtAuthenticator", slog.Any("err", err))
	}

	// middleware
	r.Use(appMiddleware.RequestSlogger)

	// handlers
	httpHandler := handlers.NewHandler(*app.Service, jwtAuthenticator)
	r.Get("/ping", httpHandler.Ping)
	r.Post("/api/user/register", httpHandler.RegisterUser)
	r.Post("/api/user/login", httpHandler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.AuthMiddleware(jwtAuthenticator))
		r.Post("/api/user/orders", httpHandler.UploadOrder)
		r.Get("/api/user/orders", httpHandler.GetUserOrders)
		r.Get("/api/user/withdrawals", httpHandler.GetUserWithdrawals)
		r.Get("/api/user/balance", httpHandler.GetBalance)
		r.Post("/api/user/balance/withdraw", httpHandler.MakeWithdrawal)
	})
	return r
}
