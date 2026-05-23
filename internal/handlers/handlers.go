package handler

import (
	"context"
	"encoding/json"
	"homework/internal/model"
	"net/http"
	"strconv"
)

type Service interface {
	CreateUser(ctx context.Context, userReq model.UserRequest) (int, error)
	DeleteUser(ctx context.Context, id int) error
	GetUser(ctx context.Context, userID string) (model.User, error)
	UpdateUser(ctx context.Context, userReq model.UserRequest) error
	GetUsersList(ctx context.Context) ([]model.User, error)
}

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user model.UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}
	
	id, err := h.svc.CreateUser(r.Context(), user)
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
	userID := r.URL.Query().Get("id")
	if userID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	user, err := h.svc.GetUser(r.Context(), userID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "can`t fiend user by id")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"User": user,
	})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	err = h.svc.DeleteUser(r.Context(), id)
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
	var user model.UserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}
	
	err = h.svc.UpdateUser(r.Context(), user)
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
	users, err := h.svc.GetUsersList(r.Context())
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
