package app

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Traliaa/KineticVPN-Bot/internal/adapter/telegram"
	"github.com/Traliaa/KineticVPN-Bot/internal/config"
	"github.com/Traliaa/KineticVPN-Bot/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type App struct {
	config   *config.Config
	bot      telegram.Bot
	Handlers *chi.Mux
	server   *http.Server
}

func NewApp() *App {

	return &App{
		config:   config.NewConfig(),
		Handlers: chi.NewRouter(),
		server:   &http.Server{},
	}

}

func (a *App) SetBot(t telegram.Bot) {
	a.bot = t
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) Start(ctx context.Context) error {
	go a.bot.Start(ctx)

	a.server.Handler = chainMiddleware(a.Handlers)

	// должен слушать все интерфейсы, поэтому не можем указать localhost
	address := fmt.Sprintf(":%d", a.config.Service.HTTPPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	logger.Info("Listening on " + address)

	return a.server.Serve(listener)
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
