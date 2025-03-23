package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testValidate/internal/auth"
)

func ProtectPersonPage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserID)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	response := map[string]interface{}{
		"message": fmt.Sprintf("Доступ разрешен, UserID: %v", userID),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProtectGreetingPage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserID)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	response := map[string]interface{}{
		"message": fmt.Sprintf("Доступ разрешен, UserID: %v", userID),
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
