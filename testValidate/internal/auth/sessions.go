package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

var sessions = make(map[string]Session)

type Session struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func CreateSession(w http.ResponseWriter, userID uuid.UUID) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24)

	session := Session{
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	sessions[sessionID] = session

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		//Secure:   true, // Рекомендуется использовать HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  expiresAt,
	}

	http.SetCookie(w, cookie)

	return sessionID, nil
}
func IsValidSession(sessionID string) bool {
	session, ok := sessions[sessionID]
	if !ok {
		return false
	}
	return session.ExpiresAt.After(time.Now())
}
func DeleteSession(sessionID string) {
	delete(sessions, sessionID)

}
func DeleteCookie(w http.ResponseWriter, cookieId string) {
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
}
