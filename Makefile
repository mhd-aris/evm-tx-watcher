run:
	go run cmd/server/main.go

lint:
	golangci-lint run ./...

build:
	go build -o bin/server cmd/server/main.go
