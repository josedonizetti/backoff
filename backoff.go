package backoff

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
	"net"
	"time"
)

var (
	// ErrAttemptIsZero is the error returned when an attempt is zero.
	ErrAttemptIsZero = errors.New("Attempt cannot be zero")
	// ErrExponentIsZero is the error returned when exponent is zero.
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

// New return a Backoff that can retries request.
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

// Get retries a request based on Backoff attempts, increasing the timeout by
// the Backoff exponent on each run.
func (b *backoff) Get(ctx context.Context, target string) (*http.Response, error) {
	var (
		resp    *http.Response
		req     *http.Request
		err     error
		timeout int
	)

	timeout = 1

	for i := 0; i < b.attempts; i++ {
		client := &http.Client{
			Timeout: time.Second * time.Duration(timeout),
		}

		req, err = http.NewRequest("GET", target, nil)
		if err != nil {
			return nil, err
		}

		req = req.WithContext(ctx)
		resp, err = client.Do(req)

		if TimeoutError(err) {
			level.Info(b.logger).Log("msg", "Request timeout", "target", target, "attempt", i+1, "timeout", timeout)

			// calculate next timeout
			timeout = timeout * b.exponent

			continue
		}

		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	return nil, err
}

// TimeoutError verifies if an error was a timeout.
func TimeoutError(err error) bool {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}
