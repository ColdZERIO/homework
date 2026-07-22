package handler

import (
	"fmt"
	"net/http"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "User ID: %s", userID)
}
