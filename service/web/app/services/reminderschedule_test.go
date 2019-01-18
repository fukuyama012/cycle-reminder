package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// フィールド更新
// modelsのテストだがトランザクション利用するのでservicesにて実施
func TestReminderSchedule_Updates(t *testing.T) {
	prepareTestDB()
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
		err := services.Transact(models.DB, func(tx *gorm.DB) error {
			rSet := models.ReminderSetting{}
			if err := rSet.GetById(tx, tt.ReminderSettingID); err != nil {
				t.Error(err)
			}

			rSch := models.ReminderSchedule{}
			if err := rSch.GetByReminderSetting(tx, rSet); err != nil {
				t.Error(err)
			}
			err := rSch.Updates(tx, tt.NotifyDate)
			assert.Nil(t, err)
			assert.Equal(t, tt.NotifyDateString, rSch.NotifyDate.Format("2006-01-02"))
			return err
		})
		assert.Nil(t, err)
	}
}

// フィールド更新 エラー
func TestReminderSchedule_UpdatesError(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		NotifyDate time.Time
		NotifyDateString string
	}{
		{time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), "2018-01-01"},
		{time.Date(2018, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), "2018-12-31"},
		{time.Date(9999, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), "9999-12-31"},
	}
	for _, tt := range tests {
		err := services.Transact(models.DB, func(tx *gorm.DB) error {
			// 空情報のままレシーバーコール
			rSch := models.ReminderSchedule{}
			err := rSch.Updates(tx, tt.NotifyDate)
			assert.Error(t, err)
			assert.Equal(t, tt.NotifyDateString, rSch.NotifyDate.Format("2006-01-02"))
			assert.Equal(t, uint(0), rSch.ID)
			return err
		})
		assert.Error(t, err)
	}
}

// 削除(チェックの整合性の為トランザクション化)
// modelsのテストだがトランザクション利用するのでservicesにて実施
func TestReminderSchedule_DeleteByReminderSettingTransaction(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  uint
	}{
		{1},
		{2},
		{3},
	}
	for _, tt := range tests {
		err := services.Transact(models.DB, func(tx *gorm.DB) error {
			recordCountBefore, errCount := models.CountReminderSchedule(tx)
			if errCount != nil {
				return errCount
			}
			rSet := models.ReminderSetting{}
			errSet := rSet.GetById(tx, tt.in)
			if errSet != nil {
				return  errSet
			}
			
			rSch := models.ReminderSchedule{}
			err := rSch.DeleteByReminderSetting(tx, rSet);
			assert.Nil(t, err)
			recordCountAfter, errCount := models.CountReminderSchedule(tx)
			if errCount != nil {
				return errCount
			}
			// レコードが減少している
			assert.Equal(t, recordCountBefore - 1, recordCountAfter)
			return nil
		})
		assert.Nil(t, err)
	}
}

// 削除エラー(チェックの整合性の為トランザクション化)
// modelsのテストだがトランザクション利用するのでservicesにて実施
func TestReminderSchedule_DeleteByReminderSettingErrorTransaction(t *testing.T) {
	prepareTestDB()
	err := services.Transact(models.DB, func(tx *gorm.DB) error {
		recordCountBefore, errCount := models.CountReminderSchedule(tx)
		if errCount != nil {
			return errCount
		}
		// 空情報
		rSet := models.ReminderSetting{}

		rSch := models.ReminderSchedule{}
		err := rSch.DeleteByReminderSetting(tx, rSet);
		assert.Error(t, err)
		recordCountAfter, errCount := models.CountReminderSchedule(tx)
		if errCount != nil {
			return errCount
		}
		// レコードが減少していない
		assert.Equal(t, recordCountBefore, recordCountAfter)
		return err
	})
	assert.Error(t, err)
}