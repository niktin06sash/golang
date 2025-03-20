package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"testValidate/internal/config"
	"testValidate/internal/database"
	"testValidate/internal/server"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.NewConfig()
	dbObject := database.DBObject{}
	dbEnvInterface := database.EnvObject{}
	db, err := database.ConnectToDb(dbObject, dbEnvInterface)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	server := server.NewServer(cfg, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v\n", sig)
		cancel()
	}()

	if err := server.Run(ctx); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}

	log.Println("Application exiting")
}
