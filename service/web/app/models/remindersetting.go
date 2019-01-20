package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"log"
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
func CreateReminderSetting(user User, name, notifyTitle, notifyText string, cycleDays uint) (*ReminderSetting, error)  {
	data, err := TransactAndReceiveData(DB, func(tx *gorm.DB) (i interface{}, e error) {
		// トランザクション内でnumber値を自動採番
		number, err := getReminderSettingsNextNumberForCreate(tx)
		if err != nil {
			return nil, err
		}
		i, e = createReminderSetting(tx, user, name, notifyTitle, notifyText, cycleDays, number)
		return
	})
	rSet, ok := data.(*ReminderSetting)
	if !ok {
		log.Panicf("cant cast ReminderSetting %#v\n", err)
	}
	return rSet, err
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
	rSet.ID = id
	if err := db.First(&rSet).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
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

// IDで削除
func (rs *ReminderSetting) DeleteById(db *gorm.DB, id uint) error {
	if id == 0 {
		return errors.New("empty ReminderSetting Id!")
	}
	rs.ID = id
	// 物理削除
	return db.Unscoped().Delete(&rs).Error
}

// 新規リマインダー作成
func createReminderSetting(db *gorm.DB, user User, name, notifyTitle, notifyText string, cycleDays, number uint) (*ReminderSetting, error) {
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
	if err := db.Create(&rSet).Error; err != nil {
		return nil, err
	}
	return &rSet, nil
}

// インサート用に次点のnumber値を取得
func getReminderSettingsNextNumberForCreate(db *gorm.DB) (uint, error) {
	type Result struct {
		Max uint
	}
	var result Result
	if err := db.Table("reminder_settings").Select("MAX(`number`) AS `max`").Scan(&result).Error; err != nil {
		return 0, err
	}
	return result.Max + uint(1), nil
}
