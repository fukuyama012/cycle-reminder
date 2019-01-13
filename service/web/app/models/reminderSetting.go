package models

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

type ReminderSetting struct {
	gorm.Model
	UserID uint `gorm:"not null;" validate:"required,numeric,min=1"`
	Number uint `gorm:"Type:smallint(5) unsigned;not null; "validate:"required,numeric,min=1"`
	Name string `gorm:"size:255;not null;" validate:"required,max=100"`
	NotifyTitle string `gorm:"size:255;not null;" validate:"max=100"`
	NotifyText string `gorm:"Type:text;not null;" validate:"required,max=1000"`
	CycleDays uint `gorm:"Type:smallint(5) unsigned;not null;" validate:"required,numeric,min=1,max=365"`
	ReminderSchedule ReminderSchedule
	ReminderLogs []ReminderLog
}

func (rSet *ReminderSetting) validate() error {
	return validator.New().Struct(*rSet)
}

// 新規リマインダー作成
func CreateReminderSetting(user User, name string, notifyTitle string, notifyText string, cycleDays uint) (*ReminderSetting, error) {
	number, err := getReminderSettingsNextNumberForInsert()
	if err != nil {
		return nil, err
	}
	rSet := ReminderSetting{
		UserID: user.ID,
		Number: number,
		Name: name,
		NotifyTitle: notifyTitle,
		NotifyText: notifyText,
		CycleDays: cycleDays,
	}
	// validator.v9
	if err := rSet.validate(); err != nil {
		return nil, err
	}
	if err := DB.Create(&rSet).Error; err != nil {
		return nil, err
	}
	return &rSet, nil
}

// リマインダー数カウント
func CountReminderSetting() (int, error) {
	var count int
	err := DB.Model(&ReminderSetting{}).Count(&count).Error
	return count, err
}

// インサート用にnumber値を取得
func getReminderSettingsNextNumberForInsert() (uint, error) {
	type Result struct {
		Max uint
	}
	var result Result
	if err := DB.Table("reminder_settings").Select("MAX(`number`) AS `max`").Scan(&result).Error; err != nil {
		return 0, err
	}
	return result.Max + uint(1), nil
}