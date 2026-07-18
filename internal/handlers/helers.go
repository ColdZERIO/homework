package handler

import "homework/internal/services"

func ToUserResponse(userDB PersistUserRequest) services.CreateUserInput {
	return services.CreateUserInput{
		Login:    userDB.Login,
		Password: userDB.Password,
		Name:     userDB.Name,
		Email:    userDB.Email,
	}
}

func ToUserListResponse(userDB UserListRequest) services.GetUserList {
	userList := make([]services.GetUser, len(userDB.Users))

	for _, value := range userDB.Users {
		userList = append(userList, services.GetUser{
			Login: value.Login,
			Name:  value.Name,
			Email: value.Email,
		})
	}

	return services.GetUserList{
		Users:  userList,
		Limit:  userDB.Limit,
		Offset: userDB.Offset,
	}
}
