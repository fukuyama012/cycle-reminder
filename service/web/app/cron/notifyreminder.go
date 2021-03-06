package main

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"log"
	"os"
	"time"
)

// Logger ログ
var Logger = log.New(os.Stdout, "", log.LstdFlags)


// コンテナ外部のcronからファイル指定実行する
func main()  {
	doNotifyReminder()
}

// doNotifyReminder 通知処理実行
func doNotifyReminder()  {
	scheduleID := uint(0)
	limit := 100
	sendTotal := 0
	sendCountOneRound := 0
	for {
		// 本日を起点に予定されている通知日時をチェック
		notifyDetails, err := services.GetRemindersReachedNotifyDate(services.GetDB(), time.Now(), scheduleID, limit)
		if err != nil {
			// 検索処理失敗
			logError("GetRemindersReachedNotifyDate %#v", err)
			break
		}
		if len(notifyDetails) == 0 {
			// 対象レコード無し
			break
		}
		scheduleID, sendCountOneRound = sendMailTargetNotifyDetails(notifyDetails)
		sendTotal += sendCountOneRound
		if scheduleID == uint(0) {
			// 一応ガード
			break
		}
	}
	logInfo("notify count %d", sendTotal)
}

// sendMailByNotifyDetails 通知詳細分メール送信実行
func sendMailTargetNotifyDetails(notifyDetails []services.NotifyDetail) (uint, int) {
	scheduleID := uint(0)
	sendCount := 0
	for _, notifyDetail := range notifyDetails {
		scheduleID = notifyDetail.ScheduleID
		// 外部API利用で問題発生をなるべく防ぐ為に余裕を持って送信する
		time.Sleep(500 * time.Millisecond)
		if err := notifyDetail.Send(); err != nil {
			logError("NotifyDetail Send ScheduleID [%d] %#v", notifyDetail.ScheduleID, err)
			continue
		}
		sendCount++
	}
	return scheduleID, sendCount
}

// logInfo 結果ログ
func logInfo(format string, v ...interface{})  {
	logfile := openFile("notify")
	defer logfile.Close()

	Logger.SetOutput(logfile)
	Logger.Printf(format, v...)
}

// logError エラーログ
func logError(format string, v ...interface{})  {
	logfile := openFile("notify_error")
	defer logfile.Close()

	Logger.SetOutput(logfile)
	Logger.Printf(format, v...)
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
