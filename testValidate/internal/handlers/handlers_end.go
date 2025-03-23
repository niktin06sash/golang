package handlers

import (
	"net/http"
	"testValidate/internal/auth"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	sessionID := sessionCookie.Value
	auth.DeleteSession(sessionID)
	newCookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
	}

	http.SetCookie(w, newCookie)
	w.WriteHeader(http.StatusOK)
}
