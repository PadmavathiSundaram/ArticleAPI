local: setup run
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