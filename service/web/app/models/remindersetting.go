package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"time"
)

type ReminderSetting struct {
	gorm.Model
	UserID uint `gorm:"not null;" validate:"required,numeric,min=1"`
	Number uint `gorm:"Type:smallint(5) unsigned;not null; "validate:"required,numeric,min=1"`
	Name string `gorm:"size:255;not null;" validate:"required,max=100"`
	NotifyTitle string `gorm:"size:255;not null;" validate:"max=100"`
	NotifyText string `gorm:"Type:text;not null;" validate:"required,max=1000"`
	CycleDays uint `gorm:"Type:smallint(5) unsigned;not null;" validate:"required,numeric,min=1,max=365"`
	//ReminderSchedule ReminderSchedule
	ReminderLogs []ReminderLog
}

func (rSet *ReminderSetting) validate() error {
	return validator.New().Struct(*rSet)
}

// 新規リマインダー作成
func CreateReminderSettingWithNumbering(userID uint, name, notifyTitle, notifyText string, cycleDays uint) (*ReminderSetting, error)  {
	data, err := TransactAndReceiveData(DB, func(tx *gorm.DB) (i interface{}, e error) {
		// トランザクション内でnumber値を自動採番
		number, err := GetReminderSettingsNextNumberForCreate(tx)
		if err != nil {
			return nil, err
		}
		i, e = CreateReminderSetting(tx, userID, name, notifyTitle, notifyText, cycleDays, number)
		return
	})
	rSet, ok := data.(*ReminderSetting)
	if !ok {
		log.Panicf("cant cast ReminderSetting %#v\n", err)
	}
	return rSet, err
}

// 新規リマインダー作成
func CreateReminderSetting(db *gorm.DB, userID uint, name, notifyTitle, notifyText string, cycleDays, number uint) (*ReminderSetting, error) {
	rSet := ReminderSetting{
		UserID: userID,
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
	if err := db.Create(&rSet).Error; err != nil {
		return nil, err
	}
	return &rSet, nil
}

// インサート用に次点のnumber値を取得
func GetReminderSettingsNextNumberForCreate(db *gorm.DB) (uint, error) {
	type Result struct {
		Max uint
	}
	var result Result
	if err := db.Table("reminder_settings").Select("MAX(`number`) AS `max`").Scan(&result).Error; err != nil {
		return 0, err
	}
	return result.Max + uint(1), nil
}

// リマインダー数カウント
func CountReminderSetting(db *gorm.DB) (int, error) {
	var count int
	err := db.Unscoped().Model(&ReminderSetting{}).Count(&count).Error
	return count, err
}

// ユーザーの全リマインド設定取得
func GetReminderSettingsByUser(db *gorm.DB, user User) ([]ReminderSetting, error) {
	var rSettings []ReminderSetting
	if err := db.Where("user_id = ?", user.ID).Find(&rSettings).Error; err != nil {
		return nil, err
	}
	return rSettings, nil
}

// IDで検索
func (rSet *ReminderSetting) GetById(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("id is 0, ReminderSetting GetById")
	}
	rSet.ID = id
	if err := db.First(&rSet).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// GetByUserIDAndNumber UserIDとNumberで検索
func (rSet *ReminderSetting) GetByUserIDAndNumber(db *gorm.DB, userID uint, number uint) error {
	if err := db.Where("user_id = ? AND number = ?", userID, number).First(&rSet).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// 起点日付＋通知間隔日数で日付を算出する 
func (rSet *ReminderSetting) CalculateNotifyDate(basisDate time.Time) time.Time {
	return basisDate.AddDate(0, 0, int(rSet.CycleDays))
}

// 更新
func (rSet *ReminderSetting) Updates(db *gorm.DB, name, notifyTitle, notifyText string, cycleDays uint) error {
	rSet.Name = name
	rSet.NotifyTitle = notifyTitle
	rSet.NotifyText = notifyText
	rSet.CycleDays = cycleDays
	// validator.v9
	if err := rSet.validate(); err != nil {
		return err
	}
	if err := db.Save(&rSet).Error; err != nil {
		return err
	}
	return nil
}

// Delete 削除
func (rSet *ReminderSetting) Delete(db *gorm.DB) error {
	if rSet.ID == 0 {
		return errors.New("empty ReminderSetting ID!")
	}
	// 物理削除
	return db.Unscoped().Delete(&rSet).Error
}



