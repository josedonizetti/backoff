package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func main() {
	exp := 1
	for i := 1; i <= 3; i++ {
		client := &http.Client{
			Timeout: time.Second * time.Duration(exp),
		}

		exp *= 2

		url := "https://httpbin.org/delay/3"
		resp, err := client.Get(url)

		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Printf("The request at %s timeout.\n", url)
			continue
		}

		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}

		fmt.Println("Request completed successfully.")
		fmt.Println(resp.StatusCode)
	}
}
