package auth

import (
	"context"
	"net/http"
)

type UserIDKey string

const UserID UserIDKey = "id"

func AuthorityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		tokenString := cookie.Value
		claims, err := Authenticate(tokenString)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			ctx := context.WithValue(r.Context(), UserID, claims["id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
func NonAuthorityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == nil {
			tokenString := cookie.Value
			_, err := Authenticate(tokenString)
			if err == nil {
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			}

		}
		next.ServeHTTP(w, r)
	})
}
