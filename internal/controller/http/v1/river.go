package v1

import (
	"context"
	"log"

	"riverqueue.com/riverui"
)

type RiverHandler struct {
	Handler *riverui.Handler
}

func NewRiverHandler(ctx context.Context, opts *riverui.HandlerOpts) RiverHandler {
	handler, err := riverui.NewHandler(opts)
	if err != nil {
		log.Fatal(err)
	}
	err = handler.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return RiverHandler{
		Handler: handler,
	}
}
