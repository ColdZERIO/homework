package services

import (
	"homework/internal/domain"
	"homework/internal/storage"
)

func ToUserResponse(userDB CreateUserInput) domain.UserOutput {
	return domain.UserOutput{
		ID:        userDB.ID,
		Login:     userDB.Login,
		Name:      userDB.Name,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
	}
}

func ToUserListResponse(userDB GetUserList) storage.UserDBListResponse {
	userList := make([]storage.UserDBResponse, len(userDB.Users))

	for _, value := range userDB.Users {
		userList = append(userList, storage.UserDBResponse{
			Login: value.Login,
			Name:  value.Name,
			Email: value.Email,
		})
	}

	return storage.UserDBListResponse{
		Users:  userList,
		Limit:  userDB.Limit,
		Offset: userDB.Offset,
	}
}
