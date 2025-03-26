FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth_service ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/auth_service /app/auth_service
COPY config.yaml /app/config.yaml
WORKDIR /app
ENTRYPOINT ["/app/auth_service"]
