test:
	go test ./... -v -cover -covermode=count -coverprofile=count.out
	go tool cover -func=count.out

lint:
	golangci-lint run

vet:
	go vet ./...


