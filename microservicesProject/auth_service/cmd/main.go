package main

import (
	"context"
	"fmt"
	"log"
	"microservicesProject/auth_service/internal/api"
	"microservicesProject/auth_service/internal/repository"
	"microservicesProject/auth_service/internal/server"
	"microservicesProject/auth_service/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
)

func main() {
	// 1. Загрузка конфигурации
	configPath := "microservicesProject/auth_service/configs/config.yml"
	repository.LoadConfig(configPath)

	// 2. Подключение к базе данных
	dbInterface := repository.DBObject{}
	db, err := repository.ConnectToDb(dbInterface)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	// 3. Инициализация слоев приложения
	repos := repository.NewAuthRepository(db)
	service := service.NewService(repos)
	handlers := api.NewHandler(service)
	srv := &server.Server{}

	// 4. Запуск сервера
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080" // Значение по умолчанию, если порт не указан в конфигурации
	}

	fmt.Printf("Starting server on port: %s\n", port) // Выводим порт для проверки

	go func() {
		if err := srv.Run(port, handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	// 5. Обработка сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Service Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("error occured on db connection close: %s", err.Error())
	}
}
