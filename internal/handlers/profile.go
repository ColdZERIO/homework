package handler

import (
	"fmt"
	"homework/internal/auth"
	"net/http"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "User ID: %s", userID)
}
