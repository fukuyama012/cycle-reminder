package models

import (
	"github.com/jinzhu/gorm"
)

type ReminderLog struct {
	gorm.Model
	UserID uint `gorm:"not null;"`
	ReminderSettingID uint `gorm:"not null;"`
	NotifyTitle string `gorm:"size:255;not null;"`
	NotifyText string `gorm:"Type:text;not null;"`
}