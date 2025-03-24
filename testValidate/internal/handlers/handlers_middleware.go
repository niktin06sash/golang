package handlers

import (
	"encoding/json"
	"net/http"
	"testValidate/internal/auth"
	"testValidate/internal/erro"

	"github.com/google/uuid"
)

type UserData struct {
	UserID uuid.UUID
}

func (handler *Handler) ProtectPersonPage(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(auth.UserID)
	if userIDValue == nil {
		http.Error(w, erro.ErrorUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	data := UserData{
		UserID: userID,
	}
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *Handler) ProtectGreetingPage(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value(auth.UserID)
	if userIDValue == nil {
		http.Error(w, erro.ErrorUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	data := UserData{
		UserID: userID,
	}
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (handler *Handler) ProtectMainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Контент для неавторизованных пользователей"))
}
