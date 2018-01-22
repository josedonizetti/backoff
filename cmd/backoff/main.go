package main

import (
	"fmt"
	"github.com/josedonizetti/backoff"
)

func main() {
	target := "https://httpbin.org/delay/3"
	attempts := 3
	exponent := 2
	request := backoff.NewRequest(attempts, exponent)
	resp, err := request.Get(target)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}

	fmt.Printf("request completed %d.\n", resp.StatusCode)
}
