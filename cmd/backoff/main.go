package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/josedonizetti/backoff"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"syscall"
)

var (
	target   = kingpin.Arg("target", "Target URL").Required().String()
	attempts = kingpin.Flag("attempts", "Number of attempts").Short('a').Default("3").Int()
	exponent = kingpin.Flag("exponent", "Tiemout exponent").Short('e').Default("2").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

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

	level.Info(logger).Log("msg", "Starting backoff")
	request, err := backoff.New(*attempts, *exponent, logger)
	if err != nil {
		level.Info(logger).Log("msg", "Error", "err", err)
		return
	}

	resp, err := request.Get(ctx, *target)

	if contextCanceled(err) {
		os.Exit(0)
	}

	if err != nil && !backoff.TimeoutError(err) {
		level.Info(logger).Log("msg", "Unexpected error", "err", err)
		os.Exit(1)
	}

	if err != nil && backoff.TimeoutError(err) {
		os.Exit(0)
	}

	level.Info(logger).Log("msg", "Request completed", "stauts_code", resp.StatusCode)
	os.Exit(0)
}

func contextCanceled(err error) bool {
	return err == context.Canceled
}
