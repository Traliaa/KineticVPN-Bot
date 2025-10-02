package client

import (
	"net/http"

	"golang.org/x/time/rate"
)

type TransportBuilder struct {
	rr http.RoundTripper
}

func NewTransportBuilder(base http.RoundTripper) *TransportBuilder {
	return &TransportBuilder{
		rr: base,
	}
}

func (b *TransportBuilder) Build() http.RoundTripper {
	return b.rr
}

func (b *TransportBuilder) WithLimiter(rl *rate.Limiter) {
	b.rr = NewLimitedTransport(b.rr, rl)
}
