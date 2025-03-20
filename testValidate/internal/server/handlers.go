package server

import (
	"encoding/json"
	"io"
	"net/http"

	"testValidate/internal/erro"
	"testValidate/internal/person"
)

func (server *Server) Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, erro.ErrorNotPost.Error(), http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	dataFromPerson, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, erro.ErrorNotReadAll.Error(), http.StatusBadRequest)
		return
	}
	var newperk = person.NewPerson()
	err = json.Unmarshal(dataFromPerson, &newperk)
	if err != nil {
		http.Error(w, erro.ErrorUnmarshal.Error(), http.StatusBadRequest)
		return
	}
	mapaerr := server.PersonService.Registration(newperk, r.Context())
	if mapaerr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": mapaerr,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Registration succesfull!",
	})

}
func (server *Server) Authentication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, erro.ErrorNotPost.Error(), http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	dataFromPerson, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, erro.ErrorNotReadAll.Error(), http.StatusBadRequest)
		return
	}
	var newperk = person.NewPerson()
	err = json.Unmarshal(dataFromPerson, &newperk)
	if err != nil {
		http.Error(w, erro.ErrorUnmarshal.Error(), http.StatusBadRequest)
		return
	}
	mapaerr := server.PersonService.Authentication(newperk, r.Context())
	if mapaerr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": mapaerr,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Authentication succesfull!",
	})

}
func (server *Server) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := server.MapaHtml["startpage"]
	tmpl.Execute(w, nil)
}
