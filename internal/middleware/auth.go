package middleware

import (
	"homework/internal/services"
	"net/http"
	"strings"
)

const (
	UserID string = "user_id"
	Login  string = "login"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "cant fiend auth header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		_, _, err := services.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "bad token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
