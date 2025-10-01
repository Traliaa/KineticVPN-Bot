package main

import (
	"context"

	"github.com/Traliaa/KineticVPN-Bot/internal/app"
)

func mustNewApp() (*app.App, context.Context) {
	ctx := context.Background()
	a := app.NewApp()

	return a, ctx
}
