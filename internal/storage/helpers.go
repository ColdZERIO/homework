package storage

import "homework/internal/domain"

func ToUser(userDB UserDB) domain.UserOutput {
	return domain.UserOutput{
		Login:     userDB.Login,
		Name:      userDB.Name,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
	}
}

func ToUserList(userListDB []UserDB) []domain.UserOutput {
	userList := make([]domain.UserOutput, len(userListDB))

	for _, user := range userListDB {
		userList = append(userList, domain.UserOutput{
			ID:        user.ID,
			Login:     user.Login,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}

	return userList
}
