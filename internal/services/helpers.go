package services

import (
	handler "homework/internal/handlers"
	"homework/internal/storage"
)

func ToUserResponse(userDB storage.UserDBResponse) handler.UserResponse {
	return handler.UserResponse{
		Login:     userDB.Login,
		Name:      userDB.Name,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
	}
}

func ToUserListResponse(userDB []storage.UserDBResponse) []handler.UserResponse {
	userList := make([]handler.UserResponse, len(userDB))

	for _, value := range userDB {
		userList = append(userList, handler.UserResponse{
			Login:     value.Login,
			Name:      value.Name,
			Email:     value.Email,
			CreatedAt: value.CreatedAt,
		})
	}

	return userList
}
