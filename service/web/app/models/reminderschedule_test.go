package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
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
		{1, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), "2018-01-08"},// 7日後
		{2, time.Date(2018, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), "2019-01-30"},// 30日後
		{3, time.Date(9998, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), "9999-01-01"},// 1日後
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

// リレーション情報で検索
func TestReminderSchedule_GetByReminderSetting(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		ReminderSettingID uint
	}{
		{1},
		{2},
	}
	for _, tt := range tests {
		rSet := models.ReminderSetting{}
		errSet := rSet.GetById(models.DB, tt.ReminderSettingID)
		assert.Nil(t, errSet)
		
		rSch := models.ReminderSchedule{}
		err := rSch.GetByReminderSetting(models.DB, rSet)
		// 存在する
		assert.Nil(t, err)
		assert.NotEqual(t, uint(0), rSch.ID)
	}
}

// リレーション情報で検索 存在しない
func TestReminderSchedule_GetByReminderSettingNotExists(t *testing.T) {
	prepareTestDB()
	// 空情報を用意
	rSet := models.ReminderSetting{}

	rSch := models.ReminderSchedule{}
	err := rSch.GetByReminderSetting(models.DB, rSet)
	// 存在しない
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.Equal(t, uint(0), rSch.ID)
}

// フィールド更新
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
		err := models.Transact(models.DB, func(tx *gorm.DB) error {
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
		err := models.Transact(models.DB, func(tx *gorm.DB) error {
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
		err := models.Transact(models.DB, func(tx *gorm.DB) error {
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
func TestReminderSchedule_DeleteByReminderSettingErrorTransaction(t *testing.T) {
	prepareTestDB()
	err := models.Transact(models.DB, func(tx *gorm.DB) error {
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

func TestCountReminderSchedule(t *testing.T) {
	count, err := models.CountReminderSchedule(models.DB)
	assert.Nil(t, err)
	assert.Equal(t, 3, count)
}