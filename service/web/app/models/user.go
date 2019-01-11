package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"size:255;not null;unique_index"`
	ReminderSettings []ReminderSetting
	ReminderLogs []ReminderLog
}

// IDで検索
func (user *User) GetById(id uint) error {
	user.ID = id
	if err := DB.First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// Eメールで検索
func (user *User) GetByEmail(email string) error {
	if err := DB.Where("email = ?", email).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// 新規ユーザー作成
func CreateUser(email string) (*User, error)  {
	user := User{
		Email: email,
	}
	if err := DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return users, nil
}
