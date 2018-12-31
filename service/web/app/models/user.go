package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email string `gorm:"size:255;not null;unique_index"`
}


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
