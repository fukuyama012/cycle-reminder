package main

import (
	"errors"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"log"
	"os"
	"time"
)

// Logger ログ
var Logger = log.New(os.Stdout, "", log.LstdFlags)


// コンテナ外部のcronからファイル指定実行する
// docker-compose exec web go run app/cron/notifyreminder.go
func main()  {
	services.InitDB()
	doNotifyReminder()
}

// doNotifyReminder 通知処理実行
func doNotifyReminder()  {
	reminderScheduleID := uint(0)
	limit := 100
	offset := 0
	sendTotal := 0
	sendCountOneRound := 0
	for {
		// 本日を起点に予定されている通知日時をチェック
		notifyDetails, err := services.GetRemindersReachedNotifyDate(services.GetDB(), time.Now(), reminderScheduleID, limit, offset)
		if err != nil {
			// 検索処理失敗
			logError("GetRemindersReachedNotifyDate %#v", err)
			break
		}
		if len(notifyDetails) == 0 {
			// 対象レコード無し
			break
		}
		reminderScheduleID, sendCountOneRound = sendMailTargetNotifyDetails(notifyDetails)
		sendTotal += sendCountOneRound
		if reminderScheduleID == uint(0) {
			// 一応ガード
			break
		}
	}
	logInfo("notify count %d", sendTotal)
}

// sendMailByNotifyDetails 通知詳細分メール送信実行
func sendMailTargetNotifyDetails(notifyDetails []services.NotifyDetail) (uint, int) {
	reminderScheduleID := uint(0)
	sendCount := 0
	for _, notifyDetail := range notifyDetails {
		reminderScheduleID = notifyDetail.ID
		if err := sendMailCore(notifyDetail); err != nil {
			logError("sendMailCore ScheduleID [%d] %#v", notifyDetail.ID, err)
			continue
		}
		// 送信成功したら本日を起点に次回通知日付更新
		if err := services.ResetReminderScheduleAfterNotify(notifyDetail.ID, time.Now()); err != nil {
			logError("ResetReminderScheduleAfterNotify ScheduleID [%d]", notifyDetail.ID)
		}
		sendCount++
		// 余裕を持って送信する
		time.Sleep(500 * time.Millisecond)
	}
	return reminderScheduleID, sendCount
}

// sendMailCore　NotifyDetailを元にメール送信し結果確認
func sendMailCore(notifyDetail services.NotifyDetail) error {
	adjustNotifyContent(&notifyDetail)
	response, err := services.SendMail(notifyDetail.Email, notifyDetail.NotifyTitle, notifyDetail.NotifyText)
	if err != nil {
		return err
	}
	if !services.IsSuccessStatusCode(response) {
		return errors.New("IsSuccessStatusCode")
	}
	return nil
}

// adjustNotifyContent　通知情報調整
func adjustNotifyContent(notifyDetail *services.NotifyDetail)  {
	// メールタイトルが空の場合、補足する
	if len(notifyDetail.NotifyTitle) == 0 {
		notifyDetail.NotifyText = "Notify from Cycle Reminder"
	}
}

// logInfo 結果ログ
func logInfo(format string, v ...interface{})  {
	logfile := openFile("notify")
	defer logfile.Close()

	Logger.SetOutput(logfile)
	Logger.Printf(format, v)
}

// logError エラーログ
func logError(format string, v ...interface{})  {
	logfile := openFile("notify_error")
	defer logfile.Close()

	Logger.SetOutput(logfile)
	Logger.Printf(format, v)
}

// openFile 指定ファイル読み込み
// （ログ基盤整備するまで暫定出力）
func openFile(fileName string) *os.File {
	logfile, err := os.OpenFile("./log/"+ fileName +".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cant open file:"+ fileName + err.Error())
	}
	return logfile
}
