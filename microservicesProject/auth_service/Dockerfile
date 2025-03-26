FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates
FROM golang:1.23.4 AS builder
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o auth_service ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates

# Копируем исполняемый файл из builder stage
COPY --from=builder /app/auth_service /usr/local/bin/auth_service

RUN chmod +x /usr/local/bin/auth_service

# Настраиваем точку входа
ENTRYPOINT ["/usr/local/bin/auth_service"]

EXPOSE 8081