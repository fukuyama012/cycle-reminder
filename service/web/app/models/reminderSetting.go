package models

import (
	"github.com/jinzhu/gorm"
)

type ReminderSetting struct {
	gorm.Model
	UserID uint `gorm:"not null;"`
	Number uint `gorm:"Type:smallint(5) unsigned;not null;"`
	Name string `gorm:"size:255;not null;"`
	NotifyTitle string `gorm:"size:255;not null;"`
	NotifyText string `gorm:"Type:text;not null;"`
	CycleDays uint `gorm:"Type:smallint(5) unsigned;not null;"`
	ReminderSchedule ReminderSchedule
	ReminderLogs []ReminderLog
}