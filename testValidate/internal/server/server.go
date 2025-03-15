package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"testValidate/internal/config"
	"testValidate/internal/database"
	"testValidate/internal/erro"
	"testValidate/internal/person"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Server struct {
	Config          *config.Config
	PersonService   *person.PersonService
	HttpServer      *http.Server
	DatabaseChecker *DatabaseChecker
}
type DatabaseChecker struct {
	DB *sql.DB
}

func NewDatabaseChecker(db *sql.DB) *DatabaseChecker {
	return &DatabaseChecker{DB: db}
}

func (c *DatabaseChecker) PingContext(ctx context.Context) error {
	return c.DB.PingContext(ctx)
}
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	validate := validator.New()
	dbRepo := database.NewDBRepository(db)
	personService := person.NewPersonService(dbRepo, validate)
	databasechecker := NewDatabaseChecker(db)
	return &Server{
		Config:          cfg,
		PersonService:   personService,
		HttpServer:      &http.Server{Addr: cfg.Port},
		DatabaseChecker: databasechecker,
	}
}

func (server *Server) Handlers() http.Handler {
	mu := mux.NewRouter()
	mu.HandleFunc("/", server.Registration).Methods("POST")
	mu.HandleFunc("/", server.GetPerson).Methods("GET")
	return mu
}
func (server *Server) Run(ctx context.Context) error {
	if err := server.DatabaseChecker.DB.PingContext(ctx); err != nil {
		log.Printf("Failed to ping database: %v", err)
		return erro.ErrorDBConnect
	}
	log.Println("Database connection verified.")
	server.HttpServer.Handler = server.Handlers()
	serverStopped := make(chan struct{})

	go func() {
		defer close(serverStopped)
		err := server.HttpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe error: %v\n", err)

		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	log.Println("Shutting down server...")
	if err := server.HttpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
		return erro.ErrorServerShutdown
	}

	<-serverStopped

	log.Println("Server stopped")
	return nil
}
