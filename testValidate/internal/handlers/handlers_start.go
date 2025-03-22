package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testValidate/internal/auth"
	"testValidate/internal/erro"
	"testValidate/internal/person"
	"time"
)

const (
	contentTypeJSON              = "application/json"
	registrationSuccessMessage   = "Registration succesfull!"
	authenticationSuccessMessage = "Authentication succesfull!"
)

func handleAuth(w http.ResponseWriter, r *http.Request, authFunc func(*person.Person, context.Context) *person.AuthenticationResponse, successMessage string) {
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
	tokenString, err := auth.GenerateJWT(authresponse.UserId.String())
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

func Registration(w http.ResponseWriter, r *http.Request, psFunc *person.PersonService) {
	handleAuth(w, r, psFunc.Registration, registrationSuccessMessage)
}

func Authentication(w http.ResponseWriter, r *http.Request, psFunc *person.PersonService) {
	handleAuth(w, r, psFunc.Authentication, authenticationSuccessMessage)
}
