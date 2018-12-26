package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email string `gorm:"size:255;not null;unique_index"`
}
