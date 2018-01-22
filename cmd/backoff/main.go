package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/josedonizetti/backoff"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-term
		level.Info(logger).Log("msg", "Stopping backoff")
		cancel()
	}()

	target := "https://httpbin.org/delay/3"
	attempts := 3
	exponent := 2

	level.Info(logger).Log("msg", "Starting backoff")
	request := backoff.NewRequest(attempts, exponent, logger)
	resp, err := request.Get(ctx, target)

	if contextCanceled(err) {
		os.Exit(0)
	}

	if err != nil {
		level.Info(logger).Log("msg", "Unexpected error", "err", err)
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "Request completed", "stauts_code", resp.StatusCode)
	os.Exit(0)
}

func contextCanceled(err error) bool {
	return err == context.Canceled
}
