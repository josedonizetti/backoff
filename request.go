package backoff

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net"
	"net/http"
	"time"
)

// Request ...
type Request struct {
	logger   log.Logger
	Attempts int
	Exponent int
}

// NewRequest ...
func NewRequest(attempts, exponent int, logger log.Logger) *Request {
	return &Request{
		Attempts: attempts,
		Exponent: exponent,
		logger:   logger,
	}
}

// Get ...
func (r *Request) Get(ctx context.Context, target string) (*http.Response, error) {
	var (
		resp    *http.Response
		err     error
		timeout int
	)

	// need to validate if attemps is 0
	// need to validate if emponent is 0

	timeout = 1

	for i := 0; i < r.Attempts; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		client := &http.Client{
			Timeout: time.Second * time.Duration(timeout),
		}

		// calculate next timeout
		timeout = timeout * r.Exponent

		resp, err = client.Get(target)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if TimeoutError(err) {
			level.Info(r.logger).Log("msg", "Request timeout", "target", target)
			continue
		}

		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	return nil, err
}

// TimeoutError ...
func TimeoutError(err error) bool {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}
