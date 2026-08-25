.PHONY: run build test vet tidy lint docker-build

build:
	go build -o bin/domain-reputation-server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -race -cover

vet:
	go vet ./...

tidy:
	go mod tidy

lint: vet test

docker-build:
	docker build -t domain-reputation-server .
