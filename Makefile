test:
	go test

clean:
	rm -rf backoff

build: clean
	go build -o backoff cmd/backoff/main.go

.PHONY: test clean build

