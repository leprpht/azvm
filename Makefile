APP_NAME=azvm

.PHONY: build run test fmt tidy clean

build:
	go build -o bin/$(APP_NAME) ./cmd/azvm

run:
	go run ./cmd/azvm

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

clean:
	rm -rf bin/