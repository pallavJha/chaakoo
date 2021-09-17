test:
	go test ./... -cover -covermode=count -coverprofile=count.out
	go tool cover -func=count.out

lint:
	golangci-lint run

vet:
	go vet ./...


