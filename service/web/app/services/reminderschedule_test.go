package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

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