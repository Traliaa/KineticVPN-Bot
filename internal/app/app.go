package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Traliaa/KineticVPN-Bot/internal/adapter/telegram"
	"github.com/Traliaa/KineticVPN-Bot/internal/config"
	"github.com/Traliaa/KineticVPN-Bot/internal/controller/middleware"
	"github.com/Traliaa/KineticVPN-Bot/pkg/logger"
	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

type App struct {
	config   *config.Config
	bot      telegram.Bot
	Handlers *chi.Mux
	Server   *http.Server
	river    *http.Server
}

func NewApp() *App {

	return &App{
		config:   config.NewConfig(),
		Handlers: chi.NewRouter(),
		Server:   &http.Server{},
		river:    &http.Server{},
	}

}

func (a *App) SetBot(t telegram.Bot) {
	a.bot = t
}
func (a *App) SetRiver(r *riverui.Handler) {
	a.river.Handler = r
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) Start(ctx context.Context) error {
	go a.bot.Start(ctx)
	a.Server.Handler = chainMiddleware(a.Handlers, middleware.Logging, middleware.TracingMiddleware, middleware.MetricsMiddleware)

	// должен слушать все интерфейсы, поэтому не можем указать localhost
	addressPublic := fmt.Sprintf(":%d", a.config.Service.PublicPort)
	addressAdmin := fmt.Sprintf(":%d", a.config.Service.AdminPort)

	listener, err := net.Listen("tcp", addressAdmin)
	if err != nil {
		return err
	}
	logger.Info("Listening on " + addressAdmin)

	go func() {

		listenerAdmin, err := net.Listen("tcp", addressPublic)
		if err != nil {
			log.Fatalf("RiverUI listening fatal: %s", err)
		}
		logger.Info("RiverUI listening on " + addressPublic)
		if err := a.river.Serve(listenerAdmin); err != nil && err != http.ErrServerClosed {
			logger.Error("riverui server failed", "err", err)
		}
	}()

	return a.Server.Serve(listener)
}

func (a *App) Stop(ctx context.Context) error {
	//go a.bot.Start(ctx)

	return nil
}

func chainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, mware := range middlewares {
		handler = mware(handler)
	}
	return handler
}
