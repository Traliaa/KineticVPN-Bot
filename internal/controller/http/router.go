package http

import (
	"context"

	"github.com/Traliaa/KineticVPN-Bot/internal/app"
	v1 "github.com/Traliaa/KineticVPN-Bot/internal/controller/http/v1"
	"riverqueue.com/riverui"
)

func AddRouter(ctx context.Context, app *app.App, opts *riverui.HandlerOpts) {
	river := v1.NewRiverHandler(ctx, opts)
	app.Handlers.Handle("/riverui/", river.Handler)
	return
}
