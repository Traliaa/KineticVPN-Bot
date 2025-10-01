package app

import (
	"context"

	"github.com/Traliaa/KineticVPN-Bot/internal/adapter/telegram"
	"github.com/Traliaa/KineticVPN-Bot/internal/config"
)

type App struct {
	config *config.Config
	bot    telegram.Bot
}

func NewApp() *App {
	return &App{
		config: config.NewConfig(),
	}

}

func (a *App) SetBot(t telegram.Bot) {
	a.bot = t
}

func (a *App) GetConfig() *config.Config {
	return a.config
}

func (a *App) Start(ctx context.Context) error {
	//go a.bot.Start(ctx)

	return nil
}
