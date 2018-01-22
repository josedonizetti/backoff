package main

import (
	"context"
	"fmt"
	"github.com/josedonizetti/backoff"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-term
		fmt.Println("Stopping...")
		cancel()
	}()

	target := "https://httpbin.org/delay/3"
	attempts := 3
	exponent := 2

	request := backoff.NewRequest(attempts, exponent)
	resp, err := request.Get(ctx, target)

	if contextCanceled(err) {
		os.Exit(0)
		return
	}

	if err != nil {
		fmt.Printf("error %v\n", err)
		os.Exit(1)
		return
	}

	fmt.Printf("request completed %d.\n", resp.StatusCode)
	os.Exit(0)
}

func contextCanceled(err error) bool {
	return err == context.Canceled
}
