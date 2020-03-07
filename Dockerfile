# build stage
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
ENV GO111MODULE=on

WORKDIR /app

# go mod setup
COPY go.mod .
COPY go.sum .
Run go mod download
RUN go mod vendor
COPY . .
# execute unit tests
RUN GO111MODULE=on CGO_ENABLED=0 go test -mod=vendor -covermode=atomic -coverprofile=coverage.out ./... 

WORKDIR /app/cmd/server

# build the app
RUN GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor -o articleapi

# final stage
FROM scratch

# Create WORKDIR
WORKDIR /app
# Copy app binary from the Builder stage image
COPY --from=builder /app/cmd/server/articleapi .
EXPOSE 4852
ENTRYPOINT ["./articleapi"]


