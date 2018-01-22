package backoff

import (
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
func (r *Request) Get(target string) (*http.Response, error) {
	var (
		resp    *http.Response
		err     error
		timeout int
	)

	// need to validate if attemps is 0
	// need to validate if emponent is 0

	timeout = 1

	for i := 1; i <= r.Attempts; i++ {
		client := &http.Client{
			Timeout: time.Second * time.Duration(timeout),
		}

		// calculate next timeout
		timeout = timeout * r.Exponent

		resp, err = client.Get(target)

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
