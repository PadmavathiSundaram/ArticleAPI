local: all run
coverage: test
all: build

build:
	docker-compose down
test:
	go test --race -v -covermode=atomic -coverprofile=coverage.out ./... 
	go tool cover -html=coverage.out
run:
	docker-compose up --build