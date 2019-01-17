package services_test

import (
	"errors"
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// インサートトランザクション成功
func TestTransactAndReceiveData(t *testing.T) {
	prepareTestDB()
	_, err := services.TransactAndReceiveData(models.DB, func(tx *gorm.DB) (i interface{}, e error) {
		number, _ := models.GetReminderSettingsNextNumberForCreate(tx)
		user := models.User{}
		errU := user.GetById(1)
		assert.Nil(t, errU)
		data, err := models.CreateReminderSetting(tx, user, "name", "title", "text", 7, number)
		assert.Nil(t, err)
		// トランザクション内でnumber増加
		numberNext, err := models.GetReminderSettingsNextNumberForCreate(tx)
		assert.Nil(t, err)
		assert.Equal(t, uint(5), numberNext)
		// errorがnilだと正常終了しCommitされる
		return data, nil
	})
	assert.Nil(t, err)
	number, err := models.GetReminderSettingsNextNumberForCreate(models.DB)
	// （トランザクション成功している）
	assert.Equal(t, uint(5), number)
	assert.Nil(t, err)
}

// インサートトランザクション失敗
func TestTransactAndReceiveDataError(t *testing.T) {
	prepareTestDB()
	_, err := services.TransactAndReceiveData(models.DB, func(tx *gorm.DB) (i interface{}, e error) {
		number, _ := models.GetReminderSettingsNextNumberForCreate(tx)
		user := models.User{}
		errU := user.GetById(1)
		assert.Nil(t, errU)
		_, err := models.CreateReminderSetting(tx, user, "name", "title", "text", 7, number)
		assert.Nil(t, err)
		// トランザクション内ではnumber増加
		numberNext, err := models.GetReminderSettingsNextNumberForCreate(tx)
		assert.Nil(t, err)
		assert.Equal(t, uint(5), numberNext)
		// errorを返却するとRollbackされる
		return nil, errors.New("tran test error")
	})
	assert.Error(t, err)
	number, err := models.GetReminderSettingsNextNumberForCreate(models.DB)
	// （トランザクション失敗）5→4へRollbackされる
	assert.Equal(t, uint(4), number)
	assert.Nil(t, err)
}

