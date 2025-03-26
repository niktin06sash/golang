FROM golang:1.21 AS builder
WORKDIR /app


COPY ./auth_service/go.mod ./auth_service/go.sum ./


RUN go mod download


COPY ./auth_service ./auth_service

RUN go build -o auth_service ./auth_service/cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/auth_service /app/auth_service
COPY ./auth_service/configs/config.yml /app/config.yaml
WORKDIR /app
ENTRYPOINT ["/app/auth_service"]
EXPOSE 8080