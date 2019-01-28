package services

import (
	"errors"
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"time"
)

// ReminderDetail リマインド詳細
type ReminderDetail struct {
	ID uint
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID uint
	Number uint
	Name string
	NotifyTitle string
	NotifyText string
	CycleDays uint
	NotifyDate time.Time
}

// NotifyDetail 通知内容詳細
type NotifyDetail struct {
	Email string
	SettingID uint
	ScheduleID uint
	NotifyTitle string
	NotifyText string
}

// CreateReminderSettingWithRelation リマインド設定と紐付くリマインド予定を作成
// basisDate  起点日付　＊基本的にはtime.Now()を指定する事になる
func CreateReminderSettingWithRelation(db *gorm.DB, userID uint, name, notifyTitle, notifyText string, cycleDays uint, basisDate time.Time) (*models.ReminderSetting, error)  {
	user := models.User{}
	// トランザクションに含めないとdeadlockする事が有る
	if err := user.GetById(db, userID); err != nil {
		return nil, err
	}
	// トランザクション内でnumber値を自動採番
	number, err := models.GetReminderSettingsNextNumberForCreate(db)
	if err != nil {
		return nil, err
	}
	rSet, errSet := models.CreateReminderSetting(db, user.ID, name, notifyTitle, notifyText, cycleDays, number)
	if errSet != nil {
		return nil, errSet
	}
	// 起点日付から登録されている間隔日数を足した日付おセット
	_, errSch := models.CreateReminderSchedule(db, rSet.ID, rSet.CalculateNotifyDate(basisDate))
	if errSch != nil {
		return nil, errSch
	}
	return rSet, nil
}

// GetReminderListByUserID ユーザー情報からリマインド一覧取得
// リマインド一覧画面等で利用
func GetReminderListByUserID(db *gorm.DB, userID uint, limit, offset int) ([]ReminderDetail, error) {
	if userID == uint(0) {
		return nil, errors.New("not exists userID, GetReminderListByUserID")
	}
	var result []ReminderDetail
	if err := db.Table("reminder_settings").Select("reminder_settings.*, reminder_schedules.notify_date").
		Joins("INNER JOIN reminder_schedules ON reminder_settings.id = reminder_schedules.reminder_setting_id").
		Where("reminder_settings.user_id = ?", userID).
		Order("reminder_settings.id", true).
		Limit(limit).Offset(offset).
		Scan(&result).Error; err != nil {
			return nil, err
	}
	return result, nil
}

// GetReminderSettingByUserIDAndNumber リマインド設定取得
func GetReminderSettingByUserIDAndNumber(db *gorm.DB, UserID, number uint) (*models.ReminderSetting, error) {
	rSet := models.ReminderSetting{}
	if err := rSet.GetByUserIDAndNumber(db, UserID, number); err != nil {
		return nil, err
	}
	return &rSet, nil
}

// UpdateReminderSettingByUserIDAndNumber リマインド設定変更
func UpdateReminderSettingByUserIDAndNumber(db *gorm.DB, userID, number uint, name, notifyTitle, notifyText string, cycleDays uint) (*models.ReminderSetting, error)  {
	data, err := models.TransactAndReceiveData(db, func(tx *gorm.DB) (interface{}, error) {
		rSet := models.ReminderSetting{}
		if err := rSet.GetByUserIDAndNumber(tx, userID, number); err != nil {
			return nil, err
		}
		if err := rSet.Updates(tx, name, notifyTitle, notifyText, cycleDays); err != nil {
			return nil, err
		}
		return &rSet, nil
	})
	if err != nil {
		return nil, err
	}
	return data.(*models.ReminderSetting), nil
}

// DeleteReminderSettingByUserIDAndNumber リマインダー設定削除
func DeleteReminderSettingByUserIDAndNumber(db *gorm.DB, userID, number uint) error {
	return models.Transact(db, func(tx *gorm.DB) error {
		rSet := models.ReminderSetting{}
		if err := rSet.GetByUserIDAndNumber(tx, userID, number); err != nil {
			return err
		}
		if err := rSet.Delete(tx); err != nil {
			return err
		}
		return nil
	})
}

// GetRemindersReachedNotifyDate 通知日付に達した全リマインド予定の通知内容取得
// メール通知処理等で利用
func GetRemindersReachedNotifyDate(db *gorm.DB, targetDate time.Time, reminderScheduleID uint, limit int) ([]NotifyDetail, error) {
	var result []NotifyDetail
	if err := db.Table("reminder_schedules").Select("reminder_settings.id AS setting_id, reminder_schedules.id AS schedule_id, users.email, reminder_settings.notify_title, reminder_settings.notify_text").
		Joins("INNER JOIN reminder_settings ON reminder_schedules.reminder_setting_id = reminder_settings.id").
		Joins("INNER JOIN users ON reminder_settings.user_id = users.id").
		Where("reminder_schedules.notify_date <= ? AND reminder_schedules.id > ?", targetDate.Format("2006-01-02"), reminderScheduleID).
		Order("reminder_schedules.id", true).
		Limit(limit).
		Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

// ResetReminderScheduleAfterNotify メール通知完了後の次回通知予定設定
// basisDate  起点日付　＊基本的にはtime.Now()を指定する事になる
func ResetReminderScheduleAfterNotify(reminderSettingID uint, basisDate time.Time) error {
	err := models.Transact(models.DB, func(tx *gorm.DB) error {
		rSet := models.ReminderSetting{}
		if err := rSet.GetById(tx, reminderSettingID); err != nil {
			return err
		}
		rSch := models.ReminderSchedule{}
		if err := rSch.GetByReminderSetting(tx, rSet); err != nil {
			return err
		}
		// 次回通知日付を起点日付から指定日数後に設定
		if err := rSch.UpdateNotifyDateDaysAfterBasis(tx, basisDate, rSet.CycleDays); err != nil {
			return err
		}
		return nil
	})
	return err
}