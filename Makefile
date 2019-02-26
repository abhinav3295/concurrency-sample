.PHONY: build

build: compile fmt vet lint

run:
	go run github.com/abhinav3295/go-meetups/cmd/concurrencysample

compile:
	go build ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint -set_exit_status ./...

test: fmt vet build
	go test ./...

test-cover:
	go test -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func=coverage.out
