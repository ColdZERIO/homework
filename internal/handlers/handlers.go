package handler

import (
	"encoding/json"
	"homework/internal/services"
	"net/http"
	"strconv"
)

// ToUser на уровень выше ??? how?

type Handler struct {
	svc services.UserService
}

func NewUserHandler(svc services.UserService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Persist(w http.ResponseWriter, r *http.Request) {
	var userReq PersistUserRequest

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	user, err := h.svc.Persist(r.Context(), userReq.ID, userReq.Login, userReq.Password, userReq.Name, userReq.Email)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant add to DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PersistUserResponse{
		ID:    user.ID,
		Login: user.Login,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (h *Handler) Find(w http.ResponseWriter, r *http.Request) {
	var user FindUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	if user.ID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	searchUser, err := h.svc.Find(r.Context(), user.ID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "can`t fiend user by id")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"User": searchUser,
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
	var req UserListRequest

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
