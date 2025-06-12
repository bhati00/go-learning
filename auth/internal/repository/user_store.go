package repository

import (
	"github.com/bhati00/go-learning/auth/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(path string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	database.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}
	DB = database
	return DB, nil
}

func SaveUser(user model.User) bool {
	if err := DB.Create(&user).Error; err != nil {
		return false
	}
	return true
}

func GetUserByUsername(username string) (model.User, bool) {
	var user model.User
	if err := DB.First(&user, "username = ?", username).Error; err != nil {
		return model.User{}, false
	}
	return user, true

}
