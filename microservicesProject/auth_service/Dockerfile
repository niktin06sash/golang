# Используем официальный образ Go в качестве базового
FROM golang:1.21 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для управления зависимостями
COPY go.mod go.sum ./

# Скачиваем и устанавливаем зависимости (используем кэширование для ускорения сборки)
RUN go mod download

# Копируем исходный код приложения
COPY . .

# Собираем приложение
RUN go build -o auth_service ./cmd/main.go

# Создаем минимальный образ для запуска приложения (используем alpine)
FROM alpine:latest

# Устанавливаем необходимые библиотеки (если требуются)
RUN apk add --no-cache ca-certificates

# Копируем собранный исполняемый файл из предыдущего этапа сборки
COPY --from=builder /app/auth_service /app/auth_service

# Копируем файл конфигурации (config.yaml)
COPY config.yaml /app/config.yaml

# Устанавливаем рабочую директорию
WORKDIR /app

# Указываем точку входа для запуска приложения
ENTRYPOINT ["/app/auth_service"]

