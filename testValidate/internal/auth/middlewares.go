package auth

import (
	"context"
	"log"
	"net/http"
	"testValidate/internal/erro"
)

type UserIDKey string

const UserID UserIDKey = "user_id"

func AuthorityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, erro.ErrorUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		sessionID := sessionCookie.Value
		if !IsValidSession(sessionID) {
			DeleteSession(sessionID)
			http.Error(w, erro.ErrorUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
		log.Println(sessionID)
		session := sessions[sessionID]
		ctx := context.WithValue(r.Context(), UserID, session.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NoAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session_id")
		if err == nil {

			sessionID := sessionCookie.Value
			if IsValidSession(sessionID) {
				w.WriteHeader(http.StatusForbidden)

				return
			} else {
				DeleteSession(sessionID)
			}
		}

		next.ServeHTTP(w, r)
	}
}
