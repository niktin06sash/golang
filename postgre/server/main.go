package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/signal"
	"postgre/database"
	"postgre/server/templates"
	"strconv"
	"syscall"
	"time"

	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Config struct {
	Enabled bool
	DB      *sql.DB
	Port    string
}
type ResponseData struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewConfig(base *sql.DB) *Config {
	return &Config{
		Enabled: true,
		DB:      base,
		Port:    ":7785",
	}
}

type Server struct {
	Config  *Config
	PS      *PersonService
	Context context.Context
	Cancel  context.CancelFunc
}

func newServer(cf *Config, newserverservice *PersonService) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		Config:  cf,
		PS:      newserverservice,
		Context: ctx,
		Cancel:  cancel,
	}
}

type PersonService struct {
	Config *Config
	Pr     *database.PersonRepository
}

func NewPersonService(config *Config, pr *database.PersonRepository) *PersonService {
	return &PersonService{
		Config: config,
		Pr:     pr,
	}
}
func (Ps *PersonService) CheckUniqueReg(v database.Users, ctx context.Context) bool {
	if Ps.Config.Enabled {
		return Ps.Pr.CheckUniqueReg(v, ctx)
	}
	return false
}
func (Ps *PersonService) InsertReg(v database.Users, ctx context.Context) (bool, int) {
	if Ps.Config.Enabled {
		return Ps.Pr.InsertReg(v, ctx)
	}
	return false, 0
}
func (ps *PersonService) CheckPassword(v database.Users, ctx context.Context) error {
	if ps.Config.Enabled {
		return ps.Pr.CheckPassword(v, ctx)
	}
	return fmt.Errorf("Сервер не включен!")
}
func (ps *PersonService) SearchPerson(v database.Users, ctx context.Context) (int, error) {
	if ps.Config.Enabled {
		return ps.Pr.SearchPerson(v, ctx)
	}
	return 0, fmt.Errorf("Сервер не включен!")
}
func (Ps *PersonService) CreateSession(v database.Users, ctx context.Context) (string, error) {
	if Ps.Config.Enabled {
		return Ps.Pr.CreateSession(v, ctx)
	}
	return "", fmt.Errorf("Сервер не включен!")
}
func (Ps *PersonService) GetSession(v database.Users, ctx context.Context) (*database.Session, error) {
	if Ps.Config.Enabled {
		return Ps.Pr.GetSession(v, ctx)
	}
	return nil, fmt.Errorf("Сервер не включен!")
}
func (Ps *PersonService) Exit(id int, cook string, ctx context.Context) error {
	if Ps.Config.Enabled {
		return Ps.Pr.Exit(id, cook, ctx)
	}
	return fmt.Errorf("Сервер не включен!")
}
func (Ps *PersonService) Delete(v database.Users, ctx context.Context) error {
	if Ps.Config.Enabled {
		return Ps.Pr.Delete(v, ctx)
	}
	return fmt.Errorf("Сервер не включен!")
}

func (Ps *PersonService) UpdateOnline(id int, flag bool, ctx context.Context) error {
	if Ps.Config.Enabled {
		return Ps.Pr.UpdateOnline(id, flag, ctx)
	}
	return fmt.Errorf("Сервер не включен!")
}
func (Ps *PersonService) GetUsers(ctx context.Context) ([]database.Users, error) {
	if Ps.Config.Enabled {
		return Ps.Pr.GetUsers(ctx)
	}
	return nil, fmt.Errorf("Сервер не включен!")
}
func (Ps *PersonService) UpdateLastEntered(id int, last time.Time, session string, ctx context.Context) error {
	if Ps.Config.Enabled {
		return Ps.Pr.UpdateLastEntered(id, last, session, ctx)
	}
	return fmt.Errorf("Сервер не включен!")
}
func (Ps *PersonService) GetLastEntered(ctx context.Context) (map[int]time.Time, error) {
	if Ps.Config.Enabled {
		return Ps.Pr.GetLastEntered(ctx)
	}
	return nil, fmt.Errorf("Сервер не включен!")
}
func (server *Server) MiddleWareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				log.Println("Отсутствует cookie сессии.")
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else {
				log.Printf("Ошибка получения cookie сессии: %v", err)
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}
func (server *Server) MainPage(we http.ResponseWriter, r *http.Request) {
	tmpl := templates.MainPage()
	tmpl.Execute(we, nil)
}
func (server *Server) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var perk = database.Users{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&perk)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	userID := perk.User_id
	nowtime := time.Now()
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println(err)
	}

	gettses := cookie.Value
	err = server.PS.UpdateLastEntered(userID, nowtime, gettses, server.Context)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := ResponseData{}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

}
func (server *Server) CheckTimeouts() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-server.Context.Done():
			return
		case <-ticker.C:
			datawithtimes, err := server.PS.GetLastEntered(server.Context)
			if err != nil {
				log.Println(err)
			}
			log.Println(datawithtimes)
			for userID, lastHeartbeatTime := range datawithtimes {

				deadline := lastHeartbeatTime.Add(20 * time.Second)
				timetocheck := time.Now()
				log.Println(timetocheck)
				if timetocheck.After(deadline) {
					err := server.PS.UpdateOnline(userID, false, server.Context)
					if err != nil {
						log.Println(err)
					}
				} else {
					err := server.PS.UpdateOnline(userID, true, server.Context)

					if err != nil {
						log.Println(err)
					}
				}
			}

		}
	}
}

func (server *Server) DataOfPerson(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	sessionID := cookie.Value
	_ = server.Config.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = $1", sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Сессия не найдена в базе данных.")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		} else {
			log.Printf("Ошибка при запросе к базе данных: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}
	}
	dataperson, err := server.PS.Pr.GetUsers(server.Context)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dataperson)
}
func (server *Server) PersonInterface(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_id")
	sessionID := cookie.Value
	var userID int
	var username string
	err := server.Config.DB.QueryRow("SELECT user_id FROM sessions WHERE session_id = $1", sessionID).Scan(&userID)
	err = server.Config.DB.QueryRow("SELECT name FROM users WHERE user_id = $1", userID).Scan(&username)
	if err != nil {

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
	tmpl := templates.PersonInterface()
	data := map[string]interface{}{
		"UserID":   userID,
		"UserName": username,
	}
	jsonData, _ := json.Marshal(data)
	err = tmpl.Execute(w, map[string]interface{}{
		"jsonData": template.JS(jsonData),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (server *Server) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
	}
	getter := r.Header.Get("Do")

	if getter == "Exit" {
		idbyte, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		id, _ := strconv.Atoi(string(idbyte))
		cookie, err := r.Cookie("session_id")
		if err != nil {
			log.Println(err)
		}

		gettses := cookie.Value
		err = server.PS.Exit(id, gettses, server.Context)
		if err != nil {
			log.Println(err)
			return
		}
		datatoreq := ResponseData{
			Success: true,
			Message: "Вы успешно вышли из аккаунта!",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(datatoreq)
	} else if getter == "Delete" {
		var newperk database.Users
		data, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(data, &newperk)
		err = server.PS.CheckPassword(newperk, server.Context)

		var datatoreq = ResponseData{}
		if err != nil {
			datatoreq = ResponseData{
				Success: false,
				Message: "Неверный пароль!",
			}
		} else {
			err := server.PS.Delete(newperk, server.Context)
			if err != nil {
				datatoreq = ResponseData{
					Success: false,
					Message: "Неверный пароль!",
				}
			} else {
				datatoreq = ResponseData{
					Success: true,
					Message: "Вы успешно удалили аккаунт!",
				}
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(datatoreq)
	}
}
func (server *Server) Authority(we http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Do") == "Registration" {
		data, _ := io.ReadAll(r.Body)
		var newperk database.Users
		err := json.Unmarshal(data, &newperk)
		newperk.Online = true
		if err != nil {
			log.Println(err)
			return
		}
		if server.PS.CheckUniqueReg(newperk, server.Context) {

			itog, id := server.PS.InsertReg(newperk, server.Context)
			newperk.User_id = id
			var datatoreq = ResponseData{}
			if !itog {
				datatoreq = ResponseData{
					Success: false,
					Message: "Данный пользователь уже зарегистрирован!",
				}
				we.Header().Set("Content-Type", "application/json")
				json.NewEncoder(we).Encode(datatoreq)
				return
			}
			sessionid, err := server.PS.CreateSession(newperk, server.Context)
			if err != nil {
				log.Println("ошибка создания сессии:", err)
				http.Error(we, "Ошибка сервера", http.StatusInternalServerError)
				return
			}
			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    sessionid,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			}
			datatoreq = ResponseData{
				Success: true,
				Message: "Регистрация успешна!",
			}
			http.SetCookie(we, cookie)
			we.Header().Set("Content-Type", "application/json")
			json.NewEncoder(we).Encode(datatoreq)

		} else {
			datatoreq := ResponseData{
				Success: false,
				Message: "Данный пользователь уже зарегистрирован!",
			}
			we.Header().Set("Content-Type", "application/json")
			json.NewEncoder(we).Encode(datatoreq)
		}

	} else if r.Header.Get("Do") == "Authentication" {
		data, _ := io.ReadAll(r.Body)
		var newperk database.Users
		err := json.Unmarshal(data, &newperk)
		if err != nil {
			log.Println(err)
			return
		}
		if !server.PS.CheckUniqueReg(newperk, server.Context) {
			var datatoreq = ResponseData{}
			if server.PS.CheckPassword(newperk, server.Context) == nil {
				id, err := server.PS.SearchPerson(newperk, server.Context)
				if err != nil {
					log.Println(err)
					return
				}
				newperk.User_id = id
				sessionid, err := server.PS.CreateSession(newperk, server.Context)
				if err != nil {
					log.Println("ошибка создания сессии:", err)
					http.Error(we, "Ошибка сервера", http.StatusInternalServerError)
					return
				}
				cookie := &http.Cookie{
					Name:     "session_id",
					Value:    sessionid,
					Path:     "/",
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				}
				http.SetCookie(we, cookie)

				datatoreq = ResponseData{
					Success: true,
					Message: "Авторизация успешна!",
				}

			} else {
				datatoreq = ResponseData{
					Success: false,
					Message: "Неверный пароль!",
				}
			}
			we.Header().Set("Content-Type", "application/json")
			json.NewEncoder(we).Encode(datatoreq)
		} else {
			datatoreq := ResponseData{
				Success: false,
				Message: "Данный пользователь не зарегистрирован!",
			}
			we.Header().Set("Content-Type", "application/json")
			json.NewEncoder(we).Encode(datatoreq)
		}
	}

}
func (server *Server) Handler() http.Handler {
	mu := mux.NewRouter()
	mu.HandleFunc("/", server.MainPage).Methods("GET")
	fhdc := http.HandlerFunc(server.Profile)
	mu.Handle("/interface", fhdc).Methods("POST")
	mu.HandleFunc("/", server.Authority).Methods("POST")
	mu.HandleFunc("/heartbeat", server.heartbeatHandler).Methods("POST")
	mu.HandleFunc("/interface", server.MiddleWareAuthentication(server.PersonInterface)).Methods("GET")
	mu.HandleFunc("/interface/persons", server.MiddleWareAuthentication(server.DataOfPerson)).Methods("GET")
	return mu
}
func (server *Server) Run() {

	httpserver := &http.Server{Addr: server.Config.Port, Handler: server.Handler()}
	go func() {
		err := httpserver.ListenAndServe()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	<-server.Context.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpserver.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

}
func main() {
	databasev := database.ConnectToDB()
	defer databasev.Close()
	config := NewConfig(databasev)
	personrepos := database.NewPersonRepository(databasev)
	personservice := NewPersonService(config, personrepos)
	server := newServer(config, personservice)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go server.CheckTimeouts()
	go func() {
		<-sigChan
		server.Cancel()
	}()
	server.Run()
}
