package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testValidate/internal/auth"
	"testValidate/internal/erro"
	"testValidate/internal/person"
)

const (
	contentTypeJSON              = "application/json"
	registrationSuccessMessage   = "Registration succesfull!"
	authenticationSuccessMessage = "Authentication succesfull!"
)

func handleAuth(w http.ResponseWriter, r *http.Request, ps person.PersonServiceInterface, successMessage string, isRegistration bool) {
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

	var newperk person.Person // Используем структуру Person
	err = json.Unmarshal(dataFromPerson, &newperk)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]interface{}{"error": erro.ErrorUnmarshal.Error()},
		})
		return
	}
	authresponse := &person.AuthenticationResponse{}
	if isRegistration {
		authresponse = ps.Registration(&newperk, r.Context())
	} else {
		authresponse = ps.Authentication(&newperk, r.Context())
	}

	if !authresponse.Success {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": authresponse.Errors,
		})
		return
	}

	_, err = auth.CreateSession(w, authresponse.UserId)
	if err != nil {
		w.Header().Set("Content-Type", contentTypeJSON)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": map[string]interface{}{"error": erro.ErrorSessionCreate.Error()},
		})
		return
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": successMessage,
	})
}

func (handler *Handler) Registration(w http.ResponseWriter, r *http.Request, ps person.PersonServiceInterface) {
	handleAuth(w, r, ps, registrationSuccessMessage, true)
}

func (handler *Handler) Authentication(w http.ResponseWriter, r *http.Request, ps person.PersonServiceInterface) {
	handleAuth(w, r, ps, authenticationSuccessMessage, false)
}
