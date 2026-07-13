package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// ToUser на уровень выше ??? how?

type UserService interface {
	Persist(ctx context.Context, userReq PersistUserRequest) (string, error)
	Delete(ctx context.Context, id int) error
	Find(ctx context.Context, userID string) (UserResponse, error)
	Update(ctx context.Context, userReq PersistUserRequest) error
	GetList(ctx context.Context, limit, offset int) (UserListResponse, error)
}

type Handler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Persist(w http.ResponseWriter, r *http.Request) {
	var user PersistUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	id, err := h.svc.Persist(r.Context(), user)
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

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	var user UserResponse

	err := json.NewDecoder(r.Body).Decode(&user)
	if user.ID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	newUser, err := h.svc.Find(r.Context(), user.ID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "can`t fiend user by id")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"User": newUser,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	var user UserResponse

	err := json.NewDecoder(r.Body).Decode(&user)
	if user.ID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	id, err := strconv.Atoi(user.ID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.svc.Delete(r.Context(), id)
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var user PersistUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	_, err = h.svc.Persist(r.Context(), user)
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

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	var req UserListResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	user, err := h.svc.GetList(r.Context(), req.Limit, req.Offset)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant get list from DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"Users": user,
	})
}
