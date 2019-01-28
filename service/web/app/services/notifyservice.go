package services

import (
	"errors"
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"time"
)

// Notify　NotifyDetailを元にメール送信し結果確認
func Notify(notifyDetail NotifyDetail) error {
	adjustNotifyContent(GetDB(), &notifyDetail)
	response, err := SendMail(notifyDetail.Email, notifyDetail.NotifyTitle, notifyDetail.NotifyText)
	if err != nil {
		return err
	}
	if !IsSuccessStatusCode(response) {
		return errors.New("IsSuccessStatusCode")
	}
	return nil
}

// adjustNotifyContent　通知情報調整
func adjustNotifyContent(db *gorm.DB, notifyDetail *NotifyDetail)  {
	// メールタイトルが空の場合、補足する
	if len(notifyDetail.NotifyTitle) == 0 {
		notifyDetail.NotifyTitle = "Notify from Cycle Reminder"
	}
	// 次回通知日付を付加する
	rSet := models.ReminderSetting{}
	if err := rSet.GetById(db, notifyDetail.SettingID); err == nil {
		notifyDetail.NotifyText += "\n\n"+ "次回通知予定: " + rSet.CalculateNotifyDate(time.Now()).Format("2006-01-02")
	}
	// フッターを付加する
	notifyDetail.NotifyText += "\n\n"+ "Cycle Reminder" +"\n"+ "http://cycle-reminder.com/"
}
