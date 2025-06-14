package server

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/dangerousmonk/gophermart/cmd/config"
	_ "github.com/dangerousmonk/gophermart/docs"
	"github.com/dangerousmonk/gophermart/internal/handlers"
	appmdwlr "github.com/dangerousmonk/gophermart/internal/middleware"
	"github.com/dangerousmonk/gophermart/internal/service"
	"github.com/dangerousmonk/gophermart/internal/utils"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
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

// Start godoc
//
//	@title						Gophermart service
//	@version					1.0
//	@description				API Server
//	@BasePath					/
//
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							Cookie
//	@name						auth
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
	r.Use(middleware.Recoverer)
	r.Use(appmdwlr.RequestSlogger)

	// handlers
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8099/swagger/doc.json"),
	))

	httpHandler := handlers.NewHandler(*app.Service, jwtAuthenticator)
	r.Get("/ping", httpHandler.Ping)
	r.Post("/api/user/register", httpHandler.RegisterUser)
	r.Post("/api/user/login", httpHandler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(appmdwlr.AuthMiddleware(jwtAuthenticator))
		r.Post("/api/user/orders", httpHandler.UploadOrder)
		r.Get("/api/user/orders", httpHandler.GetUserOrders)
		r.Get("/api/user/withdrawals", httpHandler.GetUserWithdrawals)
		r.Get("/api/user/balance", httpHandler.GetBalance)
		r.Post("/api/user/balance/withdraw", httpHandler.MakeWithdrawal)
	})
	return r
}
