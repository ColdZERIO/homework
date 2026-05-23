package model

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt int64  `json:"created_at"`
}

type UserDB struct {
	ID        int    `gorm:"column:id;primaryKey"`
	Login     string `gorm:"column:login"`
	Password  string `gorm:"column:password"`
	Name      string `gorm:"column:name"`
	Email     string `gorm:"column:email"`
	CreatedAt int64  `gorm:"column:created_at"`
}

func (UserDB) TableName() string {
	return "users"
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
