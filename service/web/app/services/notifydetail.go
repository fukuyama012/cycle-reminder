package services

import (
	"errors"
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	supplementNotifyTitle = "Send from Cycle Reminder"
	supplementNextDate    = "次回通知予定: "
	supplementSiteName    = "Cycle Reminder"
	supplementSiteURL     = "http://cycle-reminder.com/"
)

// NotifyDetail 通知内容詳細
type NotifyDetail struct {
	Email string
	SettingID uint
	ScheduleID uint
	NotifyTitle string
	NotifyText string
}

// Send メール送信し次回通知日付を更新
func (notifyDetail NotifyDetail) Send() error {
	return models.Transact(GetDB(), func(tx *gorm.DB) error {
		rSet := models.ReminderSetting{}
		if err := rSet.GetById(tx, notifyDetail.SettingID); err != nil {
			return err
		}
		if err := notifyDetail.sendCore(rSet); err != nil {
			return err
		}
		// 送信成功したら本日を起点に次回通知日付更新
		if err := ResetReminderScheduleAfterNotify(tx, rSet, time.Now()); err != nil {
			return err
		}
		return nil
	})
}

// sendCore NotifyDetailを元にメール送信し結果確認
func (notifyDetail NotifyDetail) sendCore(rSet models.ReminderSetting) error {
	notifyDetail.supplementContent(rSet)
	response, err := SendMail(notifyDetail.Email, notifyDetail.NotifyTitle, notifyDetail.NotifyText)
	if err != nil {
		return err
	}
	if !IsSuccessStatusCode(response) {
		return errors.New("IsSuccessStatusCode")
	}
	return nil
}

// supplementContent　通知情報補足
func (notifyDetail *NotifyDetail) supplementContent(rSet models.ReminderSetting)  {
	// メールタイトルが空の場合、補足する
	if len(notifyDetail.NotifyTitle) == 0 {
		notifyDetail.NotifyTitle = supplementNotifyTitle
	}
	// 次回通知日付を付加する
	notifyDetail.NotifyText += "\n\n"+ supplementNextDate + rSet.CalculateNotifyDate(time.Now()).Format("2006-01-02")
	// フッターを付加する
	notifyDetail.NotifyText += "\n\n"+ supplementSiteName +"\n"+ supplementSiteURL
}