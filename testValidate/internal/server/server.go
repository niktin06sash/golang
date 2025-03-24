package server

import (
	"context"
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"testValidate/internal/auth"
	"testValidate/internal/config"
	"testValidate/internal/database"
	"testValidate/internal/erro"
	"testValidate/internal/handlers"
	"testValidate/internal/person"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Server struct {
	Config          *config.Config
	PersonService   person.PersonServiceInterface
	HttpServer      *http.Server
	DatabaseChecker *DatabaseChecker
	MapaHtml        map[string]*template.Template
	Handler         handlers.HandlerInterface
	MuxRouter       *mux.Router
	TLSConfig       *tls.Config
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
func LoadPage(name string) (*template.Template, error) {

	tmpl, err := template.ParseFiles(name)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	mapaHtml := make(map[string]*template.Template)
	startpage, err := LoadPage("../templates/start.html")
	if err != nil {
		log.Fatal(err)
	}
	personpage, err := LoadPage("../templates/personpage.html")
	if err != nil {
		log.Fatal(err)
	}
	greetingpage, err := LoadPage("../templates/greeting.html")
	if err != nil {
		log.Fatal(err)
	}
	mapaHtml["startpage"] = startpage
	mapaHtml["personpage"] = personpage
	mapaHtml["greetingpage"] = greetingpage
	/*cert, err := loadCertsFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}*/
	validate := validator.New()
	dbRepo := database.NewDBRepository(db)
	personService := person.NewPersonService(dbRepo, validate)
	databasechecker := NewDatabaseChecker(db)
	handlers := handlers.NewHandler()
	return &Server{
		Config:          cfg,
		PersonService:   personService,
		HttpServer:      &http.Server{Addr: ":" + strconv.Itoa(cfg.Port)},
		DatabaseChecker: databasechecker,
		MapaHtml:        mapaHtml,
		Handler:         handlers,
		MuxRouter:       mux.NewRouter(),
	}
}

func (server *Server) Handlers() http.Handler {
	server.MuxRouter.HandleFunc("/api/reg", (auth.NoAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { server.Handler.Registration(w, r, server.PersonService) })))).Methods("POST")

	server.MuxRouter.HandleFunc("/api/auth", (auth.NoAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Handler.Authentication(w, r, server.PersonService)
	})))).Methods("POST")
	server.MuxRouter.HandleFunc("/profile", (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Handler.PersonPage(w, r, server.MapaHtml["personpage"])
	}))).Methods("GET")
	server.MuxRouter.HandleFunc("/", (http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Handler.MainPage(w, r, server.MapaHtml["startpage"])
	}))).Methods("GET")
	server.MuxRouter.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
		server.Handler.GreetingPage(w, r, server.MapaHtml["greetingpage"])
	}).Methods("GET")

	server.MuxRouter.Handle("/api/startpage", auth.NoAuthMiddleware(http.HandlerFunc(server.Handler.ProtectMainPage)))
	server.MuxRouter.Handle("/api/greeting", auth.AuthorityMiddleware(http.HandlerFunc(server.Handler.ProtectGreetingPage)))
	server.MuxRouter.Handle("/api/profile", auth.AuthorityMiddleware(http.HandlerFunc(server.Handler.ProtectPersonPage)))
	server.MuxRouter.Handle("/api/logout", auth.AuthorityMiddleware(http.HandlerFunc(server.Handler.LogoutHandler)))
	server.MuxRouter.Handle("/api/delete", auth.AuthorityMiddleware(http.HandlerFunc(server.Handler.DeleteHandler)))
	return server.MuxRouter
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
