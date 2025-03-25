package main

import (
	"log"
	"microservicesProject/auth_service/internal/api"
	"microservicesProject/auth_service/internal/repository"
	"microservicesProject/auth_service/internal/service"
)

func main() {
	dbObject := repository.DBObject{}
	dbEnvInterface := repository.EnvObject{}
	db, err := repository.ConnectToDb(dbObject, dbEnvInterface)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	repos := repository.NewAuthRepository(db)
	service := service.NewService(repos)
	_ = api.NewHandler(service)

}
