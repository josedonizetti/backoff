package backoff

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net"
	"net/http"
	"time"
)

var (
	// ErrAttemptIsZero ...
	ErrAttemptIsZero = errors.New("Attempt cannot be zero")
	// ErrExponentIsZero ...
	ErrExponentIsZero = errors.New("Exponent cannot be zero")
)

// Backoff ...
type Backoff interface {
	Get(context.Context, string) (*http.Response, error)
}

type backoff struct {
	logger   log.Logger
	attempts int
	exponent int
}

// New ...
func New(attempts, exponent int, logger log.Logger) (Backoff, error) {
	if attempts == 0 {
		return nil, ErrAttemptIsZero
	}

	if exponent == 0 {
		return nil, ErrExponentIsZero
	}

	return &backoff{
		attempts: attempts,
		exponent: exponent,
		logger:   logger,
	}, nil
}

// Get ...
func (b *backoff) Get(ctx context.Context, target string) (*http.Response, error) {
	var (
		resp    *http.Response
		err     error
		timeout int
	)

	timeout = 1

	for i := 0; i < b.attempts; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		client := &http.Client{
			Timeout: time.Second * time.Duration(timeout),
		}

		// calculate next timeout
		timeout = timeout * b.exponent

		resp, err = client.Get(target)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if TimeoutError(err) {
			level.Info(b.logger).Log("msg", "Request timeout", "target", target)
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
