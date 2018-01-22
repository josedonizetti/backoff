package backoff

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Request ...
type Request struct {
	Attempts int
	Exponent int
}

// NewRequest ...
func NewRequest(attempts, exponent int) *Request {
	return &Request{
		Attempts: attempts,
		Exponent: exponent,
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

		if timeoutError(err) {
			fmt.Printf("request timeout %s.\n", target)
			continue
		}

		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	return nil, err
}

func timeoutError(err error) bool {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}
