package client

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

type LimitedTransport struct {
	next http.RoundTripper
	rl   *rate.Limiter
}

func NewLimitedTransport(next http.RoundTripper, rl *rate.Limiter) *LimitedTransport {
	return &LimitedTransport{
		next: next,
		rl:   rl,
	}
}

func (l *LimitedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if err := l.rl.Wait(request.Context()); err != nil {
		return nil, fmt.Errorf("waiting limiter: %w", err)
	}

	return l.next.RoundTrip(request)
}
