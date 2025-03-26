package api

import (
	"auth_service/internal/erro"
	"auth_service/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse := BadResponse(w, erro.ErrorNotPost)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {
		jsonResponse := BadResponse(w, erro.ErrorReadAll)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		jsonResponse := BadResponse(w, erro.ErrorUnmarshal)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	response := h.services.Registrate(&newperk, r.Context())
	if !response.Success {
		jsonResponse := BadResponse(w, response.Errors.(error))
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	sessionID, time := h.services.GenerateSession(newperk.Id)
	AddCookie(w, sessionID, time)
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  response.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(jsonResponse))
}
func (h *Handler) Authorization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponse := BadResponse(w, erro.ErrorNotPost)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	datafromperson, err := io.ReadAll(r.Body)
	if err != nil {
		jsonResponse := BadResponse(w, erro.ErrorReadAll)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	var newperk model.Person
	err = json.Unmarshal(datafromperson, &newperk)
	if err != nil {
		jsonResponse := BadResponse(w, erro.ErrorUnmarshal)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	response := h.services.Authorizate(&newperk, r.Context())
	if !response.Success {
		jsonResponse := BadResponse(w, response.Errors.(error))
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	sessionID, time := h.services.GenerateSession(newperk.Id)
	AddCookie(w, sessionID, time)
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusOK)
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  response.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(jsonResponse))

}
func (h *Handler) SessionChecker(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		jsonResponse := BadResponse(w, erro.ErrorInvalidSessionID)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	response := h.services.CheckSession(sessionID)
	if !response.Success {
		jsonResponse := BadResponse(w, erro.ErrorInvalidSessionID)
		if jsonResponse != nil {
			fmt.Fprint(w, string(jsonResponse))
		}
		return
	}
	sucresponse := HTTPResponse{
		Success: true,
		UserID:  response.UserId,
	}
	jsonResponse, err := json.Marshal(sucresponse)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(jsonResponse))
}
func BadResponse(w http.ResponseWriter, err error) []byte {
	w.Header().Set("Content-Type", jsonResponseType)
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := HTTPResponse{
		Success: false,
		Errors: map[string]interface{}{
			"error": err,
		},
		UserID: uuid.Nil,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, erro.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return nil
	}
	return jsonResponse
}
func AddCookie(w http.ResponseWriter, sessionID string, duration time.Time) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		//Secure:   true, // Рекомендуется использовать HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  duration,
	}

	http.SetCookie(w, cookie)
}
