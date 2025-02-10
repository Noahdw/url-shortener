FROM golang:1.23-alpine AS base
WORKDIR /app
RUN apk add --no-cache protobuf make g++
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

FROM base AS development
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "./cmd/api"]

FROM base AS production
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api
EXPOSE 8080
CMD ["./main"]