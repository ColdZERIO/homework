package auth

import (
	"context"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyRole   contextKey = "role"
)

func CheckAccessTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("ACCESS_KEY")
		if secret == "" {
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}
		claims, err := ParseToken(parts[1], AccessToken, []byte(secret))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(ContextKeyRole) != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func CanAccessUser(r *http.Request, userID string) bool {
	callerID, ok := r.Context().Value(ContextKeyUserID).(string)
	return ok && callerID != "" && (callerID == userID || r.Context().Value(ContextKeyRole) == "admin")
}
