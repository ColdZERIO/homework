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
	service services.UserService
	tokens  *auth.TokenManager
	logger  *slog.Logger
}

func NewUserHandler(service services.UserService, tokens *auth.TokenManager, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		tokens:  tokens,
		logger:  logger.With("layer", "handler"),
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
		Role:  user.Role,
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

	pair, err := h.tokens.GeneratePair(user.ID, user.Role)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to generate token pair", slog.Any("error", err))
		jsonResponseErr(w, http.StatusInternalServerError, "cant generate tokens")
		return
	}

	h.setAuthCookies(w, pair)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"user": PersistUserResponse{
			ID:    user.ID,
			Login: user.Login,
			Role:  user.Role,
			Name:  user.Name,
			Email: user.Email,
		},
	})
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
	if !auth.CanAccessUser(r.Context(), user.ID) {
		jsonResponseErr(w, http.StatusForbidden, "access denied")
		return
	}

	searchUser, err := h.service.Find(r.Context(), user.ID)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "can`t fiend user by id")
		return
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
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body request")
		return
	}
	if user.ID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}
	if !auth.CanAccessUser(r.Context(), user.ID) {
		jsonResponseErr(w, http.StatusForbidden, "access denied")
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

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		jsonResponseErr(w, http.StatusUnauthorized, "refresh token is missing")
		return
	}

	claims, err := h.tokens.ParseRefreshToken(cookie.Value)
	if err != nil {
		jsonResponseErr(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	user, err := h.service.Find(r.Context(), claims.UserID)
	if err != nil || !user.IsActive {
		jsonResponseErr(w, http.StatusUnauthorized, "user is unavailable")
		return
	}

	pair, err := h.tokens.GeneratePair(user.ID, user.Role)
	if err != nil {
		h.logger.ErrorContext(r.Context(), "failed to refresh tokens", slog.Any("error", err))
		jsonResponseErr(w, http.StatusInternalServerError, "cant generate tokens")
		return
	}

	h.setAuthCookies(w, pair)
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, expiredCookie("access_token", "/"))
	http.SetCookie(w, expiredCookie("refresh_token", "/refresh"))
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) setAuthCookies(w http.ResponseWriter, pair auth.TokenPair) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    pair.AccessToken,
		Path:     "/",
		MaxAge:   int((15 * time.Minute).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    pair.RefreshToken,
		Path:     "/refresh",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func expiredCookie(name, path string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     path,
		MaxAge:   -1,
		HttpOnly: true,
	}
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var user PersistUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		jsonResponseErr(w, http.StatusBadRequest, "invalid body rec")
		return
	}
	if user.ID == "" {
		jsonResponseErr(w, http.StatusBadRequest, "id is required")
		return
	}
	if !auth.CanAccessUser(r.Context(), user.ID) {
		jsonResponseErr(w, http.StatusForbidden, "access denied")
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
