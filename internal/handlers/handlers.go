package handler

import (
	"encoding/json"
	"homework/internal/auth"
	"homework/internal/services"
	"log/slog"
	"net/http"
	"time"
)

type UserHandler struct {
	service services.Service
}

func NewUserHandler(service services.Service, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userReq PersistUserRequest

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	user, err := h.service.Persist(r.Context(), userReq.ID, userReq.Login, userReq.Password, userReq.Name, userReq.Email)
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var userReq PersistUserRequest

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	user, err := h.service.CheckPassword(r.Context(), userReq.Login, userReq.Password)
	if err != nil {
		jsonResponseErr(w, http.StatusUnauthorized, "invalid login or password")
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant generate access token")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant generate refresh token")
		return
	}

	// Бэк не отдает токены в куки, отдаем их через Encoder (access и refresh)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// сделать хранение в БД!!!

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refreshToken,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(7 * 24 * time.Hour),
	// 	HttpOnly: true,
	// 	Secure:   false,
	// 	SameSite: http.SameSiteLaxMode,
	// })

	err = h.service.UpdateRefreshToken(r.Context(), userReq.ID, refreshToken)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant add refresh token to database")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("authentication was successful")
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// cookie, err := r.Cookie("refresh_token")
	// if err != nil {
	// 	http.Error(w, "refresh token required", http.StatusUnauthorized)
	// 	return
	// }

	var user FindUserRequest

	token, err := h.service.GetRefreshToken(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "cant get token from DB", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ParseToken(token, auth.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	//Проверить не отозван ли токен

	accessToken, err := auth.GenerateAccessToken(claims.ID, claims.Role)
	if err != nil {
		http.Error(w, "failed to update access token", http.StatusUnauthorized)
		return
	}

	// json Encoder
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   15,
	})

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Find(w http.ResponseWriter, r *http.Request) {
	var user FindUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body request")
		return
	}

	if user.ID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	searchUser, err := h.service.Find(r.Context(), user.ID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "can`t fiend user by id")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FindUserResponse{
		ID:        searchUser.ID,
		Login:     searchUser.Login,
		Name:      searchUser.Name,
		Email:     searchUser.Email,
		CreatedAt: searchUser.CreatedAt,
	},
	)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var user DeleteUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if user.ID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}

	userDeleted, err := h.service.Delete(r.Context(), user.ID)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant delete from DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"User deleted": DeleteUserResponse{
			ID:    userDeleted.ID,
			Login: userDeleted.Login,
			Name:  userDeleted.Name,
		},
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user PersistUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	userUpdated, err := h.service.Persist(r.Context(), user.ID, user.Login, user.Password, user.Name, user.Email)
	if err != nil {
		jsonResponseErr(w, http.StatusInternalServerError, "cant update in DB")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PersistUserResponse{
		ID:    userUpdated.ID,
		Login: userUpdated.Login,
		Name:  userUpdated.Name,
		Email: userUpdated.Email,
	})
}

func (h *UserHandler) FindUserList(w http.ResponseWriter, r *http.Request) {
	var userReq FindUserListRequest

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}

	user, err := h.service.GetList(r.Context(), userReq.Limit, userReq.Offset, userReq.Where, userReq.OrderBy)
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
