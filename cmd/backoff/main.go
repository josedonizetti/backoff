package main

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/josedonizetti/backoff"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	target   = kingpin.Arg("target", "Target URL").Required().String()
	attempts = kingpin.Flag("attempts", "Number of attempts").Short('a').Default("3").Int()
	exponent = kingpin.Flag("exponent", "Timeout exponent").Short('e').Default("2").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

	if _, err := url.ParseRequestURI(*target); err != nil {
		level.Info(logger).Log("msg", "Invalid target", "err", err)
		return
	}

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
		return
	}

	if err != nil && !backoff.TimeoutError(err) {
		level.Info(logger).Log("msg", "Error", "err", err)
		return
	}

	if err != nil && backoff.TimeoutError(err) {
		return
	}

	level.Info(logger).Log("msg", "Request completed", "stauts_code", resp.StatusCode)
}

func contextCanceled(err error) bool {
	// Don't really like doing this check, but unfortunately the net/http.Client
	// wraps the context error before returning
	// https://github.com/golang/go/blob/master/src/net/http/client.go#L522
	if err, ok := err.(*url.Error); ok {
		return err.Err == context.Canceled
	}
	return false
}
