local: setup fmt lint buildMongo start
standalone: setup fmt lint run
coverage: test

setup:
	docker-compose down
	go mod download
	go mod vendor
	go mod tidy

test:
	go test --race -v -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

run:
	docker-compose up --build

buildMongo:
	docker-compose up -d mongo 
	sleep 10
	docker-compose start mongo
	
start:
	go run cmd/server/main.go -c cmd/server/config/config.local.json

# Format all go files
fmt:
	gofmt -s -w -l $(shell go list -f {{.Dir}} ./...)

# Run linters - brew install golangci/tap/golangci-lint
lint:
	golangci-lint run ./...
lint-install-mac:
	brew install golangci/tap/golangci-lint
		