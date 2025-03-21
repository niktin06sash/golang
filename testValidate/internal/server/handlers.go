package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"testValidate/internal/erro"
	"testValidate/internal/person"
)

const (
	contentTypeJSON              = "application/json"
	registrationSuccessMessage   = "Registration succesfull!"
	authenticationSuccessMessage = "Authentication succesfull!"
)

func (server *Server) handleAuth(w http.ResponseWriter, r *http.Request, authFunc func(*person.Person, context.Context) map[string]string, successMessage string) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]string{"error": erro.ErrorNotPost.Error()},
		})
		return
	}

	defer r.Body.Close()
	dataFromPerson, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]string{"error": erro.ErrorNotReadAll.Error()},
		})
		return
	}

	var newperk person.Person
	err = json.Unmarshal(dataFromPerson, &newperk)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]string{"error": erro.ErrorUnmarshal.Error()},
		})
		return
	}

	mapaerr := authFunc(&newperk, r.Context())
	if mapaerr != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": mapaerr,
		})
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": successMessage,
	})
}

func (server *Server) Registration(w http.ResponseWriter, r *http.Request) {
	server.handleAuth(w, r, server.PersonService.Registration, registrationSuccessMessage)
}

func (server *Server) Authentication(w http.ResponseWriter, r *http.Request) {
	server.handleAuth(w, r, server.PersonService.Authentication, authenticationSuccessMessage)
}
func (server *Server) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := server.MapaHtml["startpage"]
	tmpl.Execute(w, nil)
}
func (server *Server) PersonPage(w http.ResponseWriter, r *http.Request) {
	tmpl := server.MapaHtml["personpage"]
	tmpl.Execute(w, nil)
}
func (server *Server) GreetingPage(w http.ResponseWriter, r *http.Request) {
	tmpl := server.MapaHtml["greetingpage"]
	tmpl.Execute(w, nil)
}
