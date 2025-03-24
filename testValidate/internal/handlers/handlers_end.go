package handlers

import (
	"net/http"
	"testValidate/internal/auth"
)

func (handler *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	sessionID := sessionCookie.Value
	auth.DeleteSession(sessionID)
	auth.DeleteCookie(w, sessionID)
	w.WriteHeader(http.StatusOK)
}
func (handler *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {

}
