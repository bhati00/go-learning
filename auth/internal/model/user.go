package model

type User struct {
	Id       int64  `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique: not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}
