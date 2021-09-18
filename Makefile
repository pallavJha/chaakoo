test:
	go test ./... -v -cover -covermode=count -coverprofile=count.out
	go tool cover -func=count.out

lint:
	golangci-lint run

vet:
	go vet ./...

mockgen:
	mockgen -source tmux_wrapper.go -destination mocks/command_executor.go -package mocks  ICommandExecutor