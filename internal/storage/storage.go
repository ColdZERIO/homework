package storage

import (
	"context"
	handler "homework/internal/handlers"
	"homework/internal/model"
	"log"

	"gorm.io/gorm"
)

type UserDB struct {
	ID        int    `gorm:"column:id;primaryKey"`
	Login     string `gorm:"column:login"`
	Password  string `gorm:"column:password"`
	Name      string `gorm:"column:name"`
	Email     string `gorm:"column:email"`
	CreatedAt int64  `gorm:"column:created_at"`
}

type Storage struct {
	db *gorm.DB
}

func UserStorage(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Persist(ctx context.Context, userDB UserDB) (int, error) {
	err := s.db.WithContext(ctx).Create(&userDB).Error
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return userDB.ID, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	err := s.db.WithContext(ctx).Delete(&UserDB{}, id).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Storage) FindByID(ctx context.Context, id int) (model.User, error) {
	var userDB UserDB
	err := s.db.WithContext(ctx).First(&userDB, id).Error
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	return ToUser(userDB), nil
}

func (s *Storage) Update(ctx context.Context, userReq handler.UserRequest) error {
	userDB := UserDB{
		Name:      userReq.Name,
		Email:     userReq.Email,
	}

	err := s.db.WithContext(ctx).Model(&UserDB{}).Where("id = ?", userDB.ID).Updates(UserDB{Name: userDB.Name, Email: userDB.Email}).Error
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (s *Storage) GetList(ctx context.Context) ([]model.User, error) {
	var usersDB []UserDB

	err := s.db.WithContext(ctx).Find(&usersDB).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return ToUserList(usersDB), nil
}

func ToUser(userDB UserDB) model.User {
	return model.User{
		ID:        userDB.ID,
		Login:     userDB.Login,
		Name:      userDB.Name,
		Email:     userDB.Email,
		CreatedAt: userDB.CreatedAt,
	}
}

func ToUserList(userListDB []UserDB) []model.User {
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
