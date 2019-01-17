package services

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"log"
)

// 新規リマインダー作成
func CreateReminderSetting(user models.User, name, notifyTitle, notifyText string, cycleDays uint) (*models.ReminderSetting, error)  {
	data, err := TransactAndReceiveData(models.DB, func(tx *gorm.DB) (i interface{}, e error) {
		// トランザクション内でnumber値を自動採番
		number, err := models.GetReminderSettingsNextNumberForCreate(tx)
		if err != nil {
			return nil, err
		}
		i, e = models.CreateReminderSetting(tx, user, name, notifyTitle, notifyText, cycleDays, number)
		return 
	})
	rSet, ok := data.(*models.ReminderSetting)
	if !ok {
		log.Fatal("cant cast ReminderSetting")
		return nil, err
	}
	return rSet, err
}