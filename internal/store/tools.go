package store

import "homework/internal/model"

func ToUser(userDB model.UserDB) model.User {
	return model.User{
		ID:        userDB.ID,
		Login:     userDB.Login,
		Name:      userDB.Name,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
	}
}

func ToUserList(userListDB []model.UserDB) []model.User {
	userList := make([]model.User, len(userListDB))

	for _, userDB := range userListDB {
		userList = append(userList, model.User{
			ID:        userDB.ID,
			Login:     userDB.Login,
			Name:      userDB.Name,
			Email:     userDB.Email,
			CreatedAt: userDB.CreatedAt,
		})
	}

	return userList
}
