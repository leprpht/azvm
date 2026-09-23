APP_NAME=azvm

.PHONY: build run test fmt tidy clean release

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
	rm -rf bin/ dist/

release:
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build -o dist/$(APP_NAME)-v$(VERSION)-darwin-arm64 ./cmd/azvm
	GOOS=darwin GOARCH=amd64 go build -o dist/$(APP_NAME)-v$(VERSION)-darwin-amd64 ./cmd/azvm
	GOOS=linux GOARCH=amd64 go build -o dist/$(APP_NAME)-v$(VERSION)-linux-amd64 ./cmd/azvm
	GOOS=linux GOARCH=arm64 go build -o dist/$(APP_NAME)-v$(VERSION)-linux-arm64 ./cmd/azvm
	GOOS=windows GOARCH=amd64 go build -o dist/$(APP_NAME)-v$(VERSION)-windows-amd64.exe ./cmd/azvm