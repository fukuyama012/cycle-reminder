package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	gorm.Model
	Email string `gorm:"size:255;not null;unique_index" validate:"required,email"`
	ReminderSettings []ReminderSetting
	ReminderLogs []ReminderLog
}

func (user *User) validate() error {
	return validator.New().Struct(*user)
}

// 新規ユーザー作成
func CreateUser(db *gorm.DB, email string) (*User, error)  {
	user := User{
		Email: email,
	}
	// validator.v9
	if err := user.validate(); err != nil {
		return nil, err
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ユーザー数カウント
func CountUser(db *gorm.DB) (int, error) {
	var count int
	err := db.Unscoped().Model(&User{}).Count(&count).Error
	return count, err
}

// IDで検索
func (user *User) GetById(db *gorm.DB, id uint) error {
	user.ID = id
	if err := db.First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// Eメールで検索
func (user *User) GetByEmail(db *gorm.DB, email string) error {
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// ユーザー削除
func (user *User) DeleteById(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("empty userId!")
	}
	user.ID = id
	// 物理削除
	return db.Unscoped().Delete(&user).Error
}

