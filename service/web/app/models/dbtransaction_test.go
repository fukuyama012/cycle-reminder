package models_test

import (
	"errors"
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// インサートトランザクション成功
func TestTransactAndReceiveData(t *testing.T) {
	prepareTestDB()
	numberAfter, _ := models.TransactAndReceiveData(models.DB, func(tx *gorm.DB) (i interface{}, e error) {
		numberBefore, errCnt := models.CountUser(tx)
		assert.Nil(t, errCnt)
		
		_, err := models.CreateUser(tx, "trantest@example.com")
		assert.Nil(t, err)
		
		numberAfter, errCnt2 := models.CountUser(tx)
		assert.Nil(t, errCnt2)

		// トランザクション内でuser数増加
		assert.Equal(t, 1, numberAfter - numberBefore)
		// errorがnilだと正常終了しCommitされる
		return numberAfter, nil
	})
	number, err := models.CountUser(models.DB)
	assert.Nil(t, err)
	// （トランザクション成功している）
	assert.Equal(t, numberAfter, number)
}

// インサートトランザクション失敗
func TestTransactAndReceiveDataError(t *testing.T) {
	prepareTestDB()
	numberAfter, _ := models.TransactAndReceiveData(models.DB, func(tx *gorm.DB) (i interface{}, e error) {
		numberBefore, errCnt := models.CountUser(tx)
		assert.Nil(t, errCnt)

		_, err := models.CreateUser(tx, "trantest@example.com")
		assert.Nil(t, err)

		numberAfter, errCnt2 := models.CountUser(tx)
		assert.Nil(t, errCnt2)

		// トランザクション内でuser数増加
		assert.Equal(t, 1, numberAfter - numberBefore)
		// errorを返却するとRollbackされる
		return numberAfter, errors.New("tran test error")
	})
	number, err := models.CountUser(models.DB)
	assert.Nil(t, err)
	// （トランザクション失敗してRollbackしている）
	assert.NotEqual(t, numberAfter, number)
}

// インサートトランザクション成功
func TestTransact(t *testing.T) {
	prepareTestDB()
	err := models.Transact(models.DB, func(tx *gorm.DB) (e error) {
		numberBefore, errCnt := models.CountUser(tx)
		assert.Nil(t, errCnt)

		_, err := models.CreateUser(tx, "trantest@example.com")
		assert.Nil(t, err)

		numberAfter, errCnt2 := models.CountUser(tx)
		assert.Nil(t, errCnt2)

		// トランザクション内でuser数増加
		assert.Equal(t, 1, numberAfter - numberBefore)
		// errorがnilだと正常終了しCommitされる
		return nil
	})
	number, err := models.CountUser(models.DB)
	assert.Nil(t, err)
	// （トランザクション成功している）
	assert.Equal(t, 4, number)
}

// インサートトランザクション失敗
func TestTransactError(t *testing.T) {
	prepareTestDB()
	err := models.Transact(models.DB, func(tx *gorm.DB) (e error) {
		numberBefore, errCnt := models.CountUser(tx)
		assert.Nil(t, errCnt)

		_, err := models.CreateUser(tx, "trantest@example.com")
		assert.Nil(t, err)

		numberAfter, errCnt2 := models.CountUser(tx)
		assert.Nil(t, errCnt2)

		// トランザクション内でuser数増加
		assert.Equal(t, 1, numberAfter - numberBefore)
		// errorを返却するとRollbackされる
		return errors.New("tran test error")
	})
	number, err := models.CountUser(models.DB)
	assert.Nil(t, err)
	// （トランザクション失敗してRollbackしている）
	assert.Equal(t, 3, number)
}
