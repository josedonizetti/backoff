package backoff

import (
	"context"
	"github.com/go-kit/kit/log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type test struct {
	attempts         int
	expectedAttempts int
	exponent         int
	waitingTime      string
}

func TestAllAttemptsTimeoutReturnsTimeoutError(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()

	tests := []test{
		{attempts: 1, expectedAttempts: 1, exponent: 2, waitingTime: "3s"},
		{attempts: 2, expectedAttempts: 2, exponent: 2, waitingTime: "4s"},
		{attempts: 3, expectedAttempts: 3, exponent: 2, waitingTime: "5s"},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualAttempts++
			d, _ := time.ParseDuration(test.waitingTime)
			time.Sleep(d)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, _ := New(test.attempts, test.exponent, logger)
		_, err := req.Get(ctx, ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %s but got", test.expectedAttempts, actualAttempts)
			return
		}

		if err == nil {
			t.Error("expecting a timeout error, but none received")
			return
		}

		if !TimeoutError(err) {
			t.Errorf("expecting a timeout error, but got: '%v'", err)
		}
	}
}

func TestIfAnyAttemptsSucceedsReturnsResponse(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()

	tests := []test{
		{attempts: 100, expectedAttempts: 1, exponent: 2, waitingTime: "0s"},
		{attempts: 100, expectedAttempts: 2, exponent: 2, waitingTime: "1s"},
		{attempts: 100, expectedAttempts: 3, exponent: 2, waitingTime: "3s"},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualAttempts++
			d, _ := time.ParseDuration(test.waitingTime)
			time.Sleep(d)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, _ := New(test.attempts, test.exponent, logger)
		resp, err := req.Get(ctx, ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %d but got %d", test.expectedAttempts, actualAttempts)
		}

		if err != nil {
			t.Errorf("not expecting error: '%v'", err)
			return
		}

		if resp == nil {
			t.Error("not expecting nil response")
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expecting http code 200 but got %d", resp.StatusCode)
		}
	}
}

func TestExponentIncreasesOnEachAttempt(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()

	tests := []test{
		{attempts: 100, expectedAttempts: 4, exponent: 3, waitingTime: "10s"},
		{attempts: 100, expectedAttempts: 2, exponent: 15, waitingTime: "10s"},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualAttempts++
			d, _ := time.ParseDuration(test.waitingTime)
			time.Sleep(d)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req, _ := New(test.attempts, test.exponent, logger)
		resp, err := req.Get(ctx, ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %d but got %d", test.expectedAttempts, actualAttempts)
		}

		if err != nil {
			t.Errorf("not expecting error: '%v'", err)
			return
		}

		if resp == nil {
			t.Error("not expecting nil response")
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expecting http code 200 but got %d", resp.StatusCode)
		}
	}
}

func TestServerReturnDiffError(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()

	tests := []test{
		{attempts: 1, expectedAttempts: 1, exponent: 2},
		{attempts: 50, expectedAttempts: 1, exponent: 2},
		{attempts: 100, expectedAttempts: 1, exponent: 2},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			actualAttempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		req, _ := New(test.attempts, test.exponent, logger)
		resp, err := req.Get(ctx, ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %d but got %d", test.expectedAttempts, actualAttempts)
		}

		if err != nil {
			t.Errorf("not expecting error: '%v'", err)
			return
		}

		if resp == nil {
			t.Error("not expecting nil response")
			return
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expecting http code 500 but got %d", resp.StatusCode)
		}
	}
}

func TestAttemptZeroReturnError(t *testing.T) {
	logger := log.NewNopLogger()
	_, err := New(0, 2, logger)
	if err != ErrAttemptIsZero {
		t.Errorf("expecting ErrAttemptIsZero, got %v", err)
	}
}

func TestExponentZeroReturnError(t *testing.T) {
	logger := log.NewNopLogger()
	_, err := New(1, 0, logger)
	if err != ErrExponentIsZero {
		t.Errorf("expecting ErrExponentIsZero, got %v", err)
	}
}
