test:
	go test ./... -v -cover -covermode=count -coverprofile=count.out
	go tool cover -func=count.out

test-race:
	go test -race -coverprofile=coverage.txt -covermode=atomic

lint:
	golint
vet:
	go vet ./...

mockgen:
	mockgen -source tmux_wrapper.go -destination mocks/command_executor.go -package mocks  ICommandExecutor

prepare:
	go mod download
	go install golang.org/x/lint/golint@latest

build:
	$(eval version=$(shell git describe --tags --always  --abbrev=5))
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags='-extldflags=-static -w -s -X github.com/pallavJha/chaakoo/cmd.version=$(version)' -o "chaakoo-$(version)-linux-amd64" cmd/chaakoo/main.go
