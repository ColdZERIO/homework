package auth

import (
	"context"
	"net/http"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyRole   contextKey = "role"
)

func (m *TokenManager) JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := m.ParseAccessToken(cookie.Value)
		if err != nil || (claims.Role != RoleUser && claims.Role != RoleAdmin) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := RoleFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if role != RoleAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ContextKeyUserID).(string)
	return userID, ok && userID != ""
}

func RoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(ContextKeyRole).(string)
	return role, ok && role != ""
}

func CanAccessUser(ctx context.Context, requestedUserID string) bool {
	userID, userOK := UserIDFromContext(ctx)
	role, roleOK := RoleFromContext(ctx)
	return userOK && roleOK && (role == RoleAdmin || userID == requestedUserID)
}
