package model

type User struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"type:varchar(100);not null;index:name_index" json:"name"`
	Age       int    `gorm:"type:int;default 0" json:"age"`
	Email     string `gorm:"type:varchar(100);not null;default:''" json:"email"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}
