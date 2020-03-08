local: setup buildMongo start
standalone: setup run
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
	sleep 30
	docker-compose start mongo
start:
	go run cmd/server/main.go -c cmd/server/config/config.local.json
		