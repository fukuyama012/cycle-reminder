package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"time"
)

// ReminderSchedule リマインド予定
type ReminderSchedule struct {
	gorm.Model
	ReminderSettingID uint `gorm:"not null;" validate:"required,min=1"`
	NotifyDate time.Time `gorm:"Type:date;not null;" validate:"required,date"`
}

func (rSch *ReminderSchedule) validate() error {
	validate := validator.New()
	// 日付フォーマット向けにカスタムバリデーション追加
	if err := validate.RegisterValidation("date", isDateFormat); err != nil {
		return err
	}
	return validate.Struct(*rSch)
}

// カスタムバリデーションの詳細
func isDateFormat(fl validator.FieldLevel) bool {
	rSet, ok := fl.Top().Interface().(ReminderSchedule)
	if !ok {
		return false
	}
	// date部分の文字列だけ取り出して比較
	_, err := time.Parse("2006-01-02", rSet.NotifyDate.String()[0:10])
	if err != nil {
		return false
	}
	return true
}

// CreateReminderSchedule 新規リマインダー予定作成
func CreateReminderSchedule(db *gorm.DB, rSet ReminderSetting, basisDate time.Time) (*ReminderSchedule, error) {
	rSch := ReminderSchedule{
		ReminderSettingID: rSet.ID,
		NotifyDate: rSet.CalculateNotifyDate(basisDate), // 起点日付から通知間隔日数を利用し計算
	}
	// validator.v9
	if err := rSch.validate(); err != nil {
		return nil, err
	}
	if err := db.Create(&rSch).Error; err != nil {
		return nil, err
	}
	return &rSch, nil
}

// CountReminderSchedule 全数カウント
func CountReminderSchedule(db *gorm.DB) (int, error) {
	var count int
	err := db.Unscoped().Model(&ReminderSchedule{}).Count(&count).Error
	return count, err
}

// GetByReminderSetting リマインド設定（ユニークキー）で検索
func (rSch *ReminderSchedule) GetByReminderSetting(db *gorm.DB, rSet ReminderSetting) error {
	if err := db.Where("reminder_setting_id = ?", rSet.ID).First(&rSch).Error; err != nil {
		if gorm.IsRecordNotFoundError(err){
			return gorm.ErrRecordNotFound
		}
		return err
	}
	return nil
}

// Updates 更新
func (rSch *ReminderSchedule) Updates(db *gorm.DB, notifyDate time.Time) error {
	rSch.NotifyDate = notifyDate
	// validator.v9
	if err := rSch.validate(); err != nil {
		return err
	}
	if err := db.Save(&rSch).Error; err != nil {
		return err
	}
	return nil
}

// UpdateNotifyDateDaysAfterBasis 通知日付を指定日時から指定日数後に変更する
func (rSch *ReminderSchedule) UpdateNotifyDateDaysAfterBasis(db *gorm.DB, basisDate time.Time, daysAfter uint) error {
	return rSch.Updates(db, basisDate.AddDate(0, 0, int(daysAfter)))
}

// DeleteByReminderSetting リマインド設定（ユニークキー）で削除
func (rSch *ReminderSchedule) DeleteByReminderSetting(db *gorm.DB, rSet ReminderSetting) error {
	if rSet.ID == 0 {
		return errors.New("cant delete ReminderSchedule, empty reminder_setting_id")
	}
	// 物理削除
	return db.Unscoped().Delete(ReminderSchedule{}, "reminder_setting_id = ?", rSet.ID).Error
}