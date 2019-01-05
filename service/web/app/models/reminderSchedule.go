package models

import (
	"github.com/jinzhu/gorm"
)

type ReminderSchedule struct {
	gorm.Model
	ReminderSettingID uint `gorm:"not null;"`
	NotifyDate string `gorm:"Type:date;not null;"`
}