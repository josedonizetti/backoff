COMMIT_HASH ?= $(shell git rev-parse HEAD)
DOCKER_IMAGE ?= "josedonizetti/backoff:$(COMMIT_HASH)"

test:
	go test

clean:
	rm -rf backoff

build: clean
	go build -o backoff cmd/backoff/main.go

docker:
	rm -rf backoff-linux-amd64
	GOOS=linux GOARCH=amd64 go build -o backoff-linux-amd64 cmd/backoff/main.go
	docker build . -t $(DOCKER_IMAGE)
	docker tag $(DOCKER_IMAGE) josedonizetti/backoff:latest
	rm -rf backoff-linux-amd64

.PHONY: test clean build docker
