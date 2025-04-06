FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev git protoc

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p app/sdk/proto/mlog

RUN protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative app/sdk/proto/mlog/logs.proto

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go/bin/server ./cmd/server

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /app/exports

COPY --from=builder /go/bin/server /app/server

EXPOSE 50051

CMD ["/app/server"]