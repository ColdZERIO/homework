package storage

func ToUser(userDB UserDB) UserDBResponse {
	return UserDBResponse{
		ID:        userDB.ID,
		Login:     userDB.Login,
		Name:      userDB.Name,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
	}
}

func ToUserList(userListDB []UserDB) []UserDBResponse {
	userList := make([]UserDBResponse, len(userListDB))

	for _, userDB := range userListDB {
		userList = append(userList, UserDBResponse{
			ID:        userDB.ID,
			Login:     userDB.Login,
			Name:      userDB.Name,
			Email:     userDB.Email,
			CreatedAt: userDB.CreatedAt,
		})
	}

	return userList
}
