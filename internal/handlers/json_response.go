package handler

import (
	"encoding/json"
	"net/http"
)

func jsonResponseErr(w http.ResponseWriter, statusCode int, msg any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]any{
		"error": msg,
	})
}
