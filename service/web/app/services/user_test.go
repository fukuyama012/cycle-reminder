package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ユーザー削除(チェックの整合性の為トランザクション化)
// models.DeleteUserByIdのテストだがトランザクション利用するのでservicesにて実施
func TestUser_DeleteUserByIdTransaction(t *testing.T) {
	tests := []struct {
		in  uint
		out bool
	}{
		{1, true},
		{2, true},
		{9999, false},
	}
	for _, tt := range tests {
		err := services.Transact(models.DB, func(tx *gorm.DB) error {
			recordCountBefore, errCount := models.CountUser(tx)
			if errCount != nil {
				return errCount
			}
			user := models.User{}
			err := user.DeleteById(tx, tt.in);
			if err != nil {
				return err
			}
			recordCountAfter, errCount := models.CountUser(tx)
			if errCount != nil {
				return errCount
			}
			if tt.out {
				// レコードが減少している
				assert.Equal(t, recordCountBefore - 1, recordCountAfter)
			} else {
				// 存在しないユーザー
				// レコードが減少していない
				assert.Equal(t, recordCountBefore, recordCountAfter)
			}
			return nil
		})
		assert.Nil(t, err)
	}
}

// ユーザー削除エラー(チェックの整合性の為トランザクション化)
// models.DeleteUserByIdのテストだがトランザクション利用するのでservicesにて実施
func TestUser_DeleteUserByIdErrorTransaction(t *testing.T) {
	tests := []struct {
		in  uint
	}{
		{0}, // ID指定が不正
	}
	for _, tt := range tests {
		err := services.Transact(models.DB, func(tx *gorm.DB) error {
			recordCountBefore, errCount := models.CountUser(tx)
			if errCount != nil {
				return errCount
			}
			user := models.User{}
			err := user.DeleteById(tx, tt.in);
			assert.Error(t, err)
			recordCountAfter, errCount := models.CountUser(tx)
			if errCount != nil {
				return errCount
			}
			// id=0指定エラー時
			// レコードが減少していない
			assert.Equal(t, recordCountBefore, recordCountAfter)
			return err
		})
		assert.Error(t, err)
	}
}