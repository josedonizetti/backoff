package backoff

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// tests
type test struct {
	attempts         int
	expectedAttempts int
	exponent         int
	waitingTime      string
}

func TestAllAttemptsTimeoutReturnsTimeoutError(t *testing.T) {
	tests := []test{
		{attempts: 1, expectedAttempts: 1, exponent: 2, waitingTime: "3s"},
		{attempts: 2, expectedAttempts: 2, exponent: 2, waitingTime: "4s"},
		{attempts: 3, expectedAttempts: 3, exponent: 2, waitingTime: "5s"},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: this should probably be atomic
			actualAttempts++
			d, _ := time.ParseDuration(test.waitingTime)
			time.Sleep(d)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req := NewRequest(test.attempts, test.exponent)
		_, err := req.Get(ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %s but got", test.expectedAttempts, actualAttempts)
			return
		}

		if err == nil {
			t.Error("expecting a timeout error, but none received")
			return
		}

		if !timeoutError(err) {
			t.Errorf("expecting a timeout error, but got: '%v'", err)
		}
	}
}

func TestIfAnyAttemptsSucceedsReturnsResponse(t *testing.T) {
	tests := []test{
		{attempts: 100, expectedAttempts: 1, exponent: 2, waitingTime: "0s"},
		{attempts: 100, expectedAttempts: 2, exponent: 2, waitingTime: "1s"},
		{attempts: 100, expectedAttempts: 3, exponent: 2, waitingTime: "3s"},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: this should probably be atomic
			actualAttempts++
			d, _ := time.ParseDuration(test.waitingTime)
			time.Sleep(d)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req := NewRequest(test.attempts, test.exponent)
		resp, err := req.Get(ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %d but got %d", test.expectedAttempts, actualAttempts)
		}

		if err != nil {
			t.Errorf("not expected error: '%v'", err)
			return
		}

		if resp == nil {
			t.Error("not expected nil response")
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected http code 200 but got %d", resp.StatusCode)
		}
	}
}

func TestExponentIncreasesOnEachAttempt(t *testing.T) {
	tests := []test{
		{attempts: 100, expectedAttempts: 4, exponent: 3, waitingTime: "10s"},
		{attempts: 100, expectedAttempts: 2, exponent: 15, waitingTime: "10s"},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: this should probably be atomic
			actualAttempts++
			d, _ := time.ParseDuration(test.waitingTime)
			time.Sleep(d)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		req := NewRequest(test.attempts, test.exponent)
		resp, err := req.Get(ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %d but got %d", test.expectedAttempts, actualAttempts)
		}

		if err != nil {
			t.Errorf("not expected error: '%v'", err)
			return
		}

		if resp == nil {
			t.Error("not expected nil response")
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected http code 200 but got %d", resp.StatusCode)
		}
	}
}

// diff error
// 3 attempts and diff error
// 2 attempts and diff error
// 1 attempt and diff error
func TestServerReturnDiffError(t *testing.T) {
	tests := []test{
		{attempts: 1, expectedAttempts: 1, exponent: 2},
		{attempts: 50, expectedAttempts: 1, exponent: 2},
		{attempts: 100, expectedAttempts: 1, exponent: 2},
	}

	for _, test := range tests {
		actualAttempts := 0
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: this should probably be atomic
			actualAttempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		req := NewRequest(test.attempts, test.exponent)
		resp, err := req.Get(ts.URL)

		if test.expectedAttempts != actualAttempts {
			t.Errorf("expecting %d but got %d", test.expectedAttempts, actualAttempts)
		}

		if err != nil {
			t.Errorf("not expected error: '%v'", err)
			return
		}

		if resp == nil {
			t.Error("not expected nil response")
			return
		}

		if resp.StatusCode != http.StatusInternalServerError {
			t.Errorf("expected http code 500 but got %d", resp.StatusCode)
		}
	}
}
