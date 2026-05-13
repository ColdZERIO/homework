package handler

import (
	"context"
	"encoding/json"
	"homework/internal/model"
	"net/http"
	"strconv"
)

type Service interface {
	CreateUser(ctx context.Context, user model.User) (int, error)
	DeleteUser(ctx context.Context, id int) error
	GetUser(ctx context.Context, id int) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	GetUsersList(ctx context.Context) ([]model.User, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponseErr(w, http.StatusMethodNotAllowed, "method not allowed (GET only)")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonResponseErr(w, http.StatusMethodNotAllowed, "method not allowed (POST only)")
		return
	}

	var user model.User
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}
	
	id, err := h.svc.CreateUser(ctx, user)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant add to DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"UserID": id,
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponseErr(w, http.StatusMethodNotAllowed, "method not allowed (GET only)")
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.svc.GetUser(context.Background(), id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"User": user,
	})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		jsonResponseErr(w, http.StatusMethodNotAllowed, "method not allowed (DELETE only)")
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.svc.DeleteUser(context.Background(), id)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant delete from DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "user deleted",
	})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		jsonResponseErr(w, http.StatusMethodNotAllowed, "method not allowed (PUT only)")
		return
	}

	var user model.User
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}
	
	err = h.svc.UpdateUser(ctx, user)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant update in DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"User": user,
	})
}

func (h *Handler) GetUsersList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		jsonResponseErr(w, http.StatusMethodNotAllowed, "method not allowed (GET only)")
		return
	}

	users, err := h.svc.GetUsersList(context.Background())
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant get from DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"Users": users,
	})
}
