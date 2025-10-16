run:
	go run ./cmd/server

tidy:
	go mod tidy

test:
	go test ./...

lint:
	golangci-lint run || true
