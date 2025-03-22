package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"testValidate/internal/erro"
	"testValidate/internal/person"
)

const (
	contentTypeJSON              = "application/json"
	registrationSuccessMessage   = "Registration succesfull!"
	authenticationSuccessMessage = "Authentication succesfull!"
)

func (server *Server) handleAuth(w http.ResponseWriter, r *http.Request, authFunc func(*person.Person, context.Context) *person.AuthenticationResponse, successMessage string) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]interface{}{"error": erro.ErrorNotPost.Error()},
		})
		return
	}

	defer r.Body.Close()
	dataFromPerson, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]interface{}{"error": erro.ErrorNotReadAll.Error()},
		})
		return
	}

	var newperk person.Person
	err = json.Unmarshal(dataFromPerson, &newperk)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]interface{}{"error": erro.ErrorUnmarshal.Error()},
		})
		return
	}

	authresponse := authFunc(&newperk, r.Context())
	if !authresponse.Success {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": authresponse.Errors,
		})
		return
	}
	tokenString, err := GenerateJWT(authresponse.UserId.String())
	if err != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]interface{}{"error": erro.ErrorJWTCreate.Error()},
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
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
