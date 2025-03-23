package auth

import (
	"time"

	"github.com/google/uuid"
)

var sessions = make(map[string]Session)

type Session struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func CreateSession(userID uuid.UUID) (string, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24)
	sessions[sessionID] = Session{UserID: userID, ExpiresAt: expiresAt}
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
