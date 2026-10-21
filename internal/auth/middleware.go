package auth

import (
	"context"
	"encoding/json"
	"homework/internal/domain"
	"net/http"
	"os"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
)

func CheckAccessTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := ParseToken(cookie.Value, AccessToken, []byte(os.Getenv("ACCESS_TOKEN")))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Проверка и сравнение ID и Роли

		var user domain.UserRequest

		err = json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "cant read json", http.StatusBadRequest)
			return
		}

		if claims.ID != user.ID && claims.Role != user.Role {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
