package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// 新規作成
func TestCreateReminderSchedule(t *testing.T) {
	prepareTestDB()
	// ユニークキー制約でインサート出来ないのでレコード全削除しておく
	if err := models.DB.Unscoped().Delete(&models.ReminderSchedule{}).Error; err != nil {
		t.Error(err)
	}
	tests := []struct {
		ReminderSettingID  uint
		NotifyDate time.Time
		NotifyDateString string
	}{
		{1, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), "2018-01-01"},
		{2, time.Date(2018, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), "2018-12-31"},
		{3, time.Date(9999, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), "9999-12-31"},
	}
	for _, tt := range tests {
		rSet := models.ReminderSetting{}
		if err := rSet.GetById(models.DB, tt.ReminderSettingID); err != nil {
			t.Error(err)
		}
		rSch, err := models.CreateReminderSchedule(models.DB, rSet, tt.NotifyDate)
		assert.Nil(t, err)
		// 日付が正常に設定されている
		assert.Equal(t, tt.NotifyDateString, rSch.NotifyDate.Format("2006-01-02"))
		assert.NotEqual(t, uint(0), rSch.ReminderSettingID)
	}
}

// 新規作成 バリデーションエラー
func TestCreateReminderScheduleError(t *testing.T) {
	prepareTestDB()
	// ユニークキー制約でインサート出来ないのでレコード全削除しておく
	if err := models.DB.Unscoped().Delete(&models.ReminderSchedule{}).Error; err != nil {
		t.Error(err)
	}
	rSetEmpty := models.ReminderSetting{}
	tests := []struct {
		rSet models.ReminderSetting
		NotifyDate time.Time
	}{
		{rSetEmpty, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())}, // ReminderSettingIDが不正
	}
	for _, tt := range tests {
		rSch, err := models.CreateReminderSchedule(models.DB, tt.rSet, tt.NotifyDate)
		assert.Error(t, err)
		// 正常に設定されていない
		assert.Nil(t, rSch)
	}
}
