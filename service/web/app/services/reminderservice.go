package services

import (
	"errors"
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// ReminderList リマインド詳細
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
	ReminderSettingID uint
	NotifyTitle string
	NotifyText string
}

// CreateReminderSettingWithRelation リマインド設定と紐付くリマインド予定を作成
// basisDate  起点日付　＊基本的にはtime.Now()を指定する事になる
func CreateReminderSettingWithRelation(userID uint, name, notifyTitle, notifyText string, cycleDays uint, basisDate time.Time) (*models.ReminderSetting, error)  {
	data, err := models.TransactAndReceiveData(models.DB, func(tx *gorm.DB) (interface{}, error) {
		user := models.User{}
		// トランザクションに含めないとdeadlockする事が有る
		if err := user.GetById(tx, userID); err != nil {
			return nil, err
		}
		// トランザクション内でnumber値を自動採番
		number, err := models.GetReminderSettingsNextNumberForCreate(tx)
		if err != nil {
			return nil, err
		}
		rSet, errSet := models.CreateReminderSetting(tx, user.ID, name, notifyTitle, notifyText, cycleDays, number)
		if errSet != nil {
			return nil, errSet
		}
		// 起点日付から登録されている間隔日数を足した日付おセット
		_, errSch := models.CreateReminderSchedule(tx, rSet.ID, rSet.CalculateNotifyDate(basisDate))
		if errSch != nil {
			return nil, errSch
		}
		return rSet, nil
	})
	if err != nil {
		return nil, err
	}
	rSet, ok := data.(*models.ReminderSetting)
	if !ok {
		log.Panicf("cant cast ReminderSetting %#v\n", data)
	}
	return rSet, nil
}

// GetReminderListByUser ユーザー情報からリマインド一覧取得
func GetReminderListByUser(user models.User, limit, offset int) ([]ReminderDetail, error) {
	if user.ID == uint(0) {
		return nil, errors.New("not exists userID, GetReminderListByUser")
	}
	var result []ReminderDetail
	if err := models.DB.Table("reminder_settings").Select("reminder_settings.*, reminder_schedules.notify_date").
		Joins("LEFT JOIN reminder_schedules ON reminder_settings.id = reminder_schedules.reminder_setting_id").
		Where("user_id = ?", user.ID).
		Limit(limit).Offset(offset).
		Scan(&result).Error; err != nil {
			return nil, err
	}
	return result, nil
}

// GetReminderSchedulesReachedNotifyDate 通知日付に達した全リマインド予定の通知内容取得
func GetRemindersReachedNotifyDate(db *gorm.DB, targetDate time.Time, limit, offset int) ([]NotifyDetail, error) {
	var result []NotifyDetail
	if err := db.Table("users").Select("users.email, reminder_schedules.reminder_setting_id, reminder_settings.notify_title, reminder_settings.notify_text").
		Joins("LEFT JOIN reminder_settings ON users.id = reminder_settings.user_id").
		Joins("LEFT JOIN reminder_schedules ON reminder_settings.id = reminder_schedules.reminder_setting_id").
		Where("reminder_schedules.notify_date <= ?", targetDate.Format("2006-01-02")).
		Limit(limit).Offset(offset).
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