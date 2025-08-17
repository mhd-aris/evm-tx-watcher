run:
	go run cmd/api/main.go

lint:
	golangci-lint run ./...

build:
	go build -o bin/api cmd/api/main.go
