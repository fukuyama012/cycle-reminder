package services

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

// CreateReminderSettingWithRelation リマインド設定と紐付くリマインド予定を作成
func CreateReminderSettingWithRelation(user models.User, name, notifyTitle, notifyText string, cycleDays uint, basisDate time.Time) (*models.ReminderSetting, error)  {
	data, err := models.TransactAndReceiveData(models.DB, func(tx *gorm.DB) (interface{}, error) {
		// トランザクション内でnumber値を自動採番
		number, err := models.GetReminderSettingsNextNumberForCreate(tx)
		if err != nil {
			return nil, err
		}
		rSet, errSet := models.CreateReminderSetting(tx, user, name, notifyTitle, notifyText, cycleDays, number)
		if errSet != nil {
			return nil, errSet
		}
		_, errSch := models.CreateReminderSchedule(tx, *rSet, basisDate)
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
		log.Panicf("cant cast ReminderSetting %#v\n", err)
	}
	return rSet, err
}
