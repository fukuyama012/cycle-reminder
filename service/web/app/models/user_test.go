package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 新規ユーザー作成
func TestUserCreateUser(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  string
		out string
	}{
		{"insert1@example.com", "insert1@example.com"},
		{"insert2@example.com", "insert2@example.com"},
	}
	for _, tt := range tests {
		user, err := models.CreateUser(models.DB, tt.in)
		if err != nil {
			t.Error(err)
		}
		// メールアドレスが正常に設定されている
		assert.Equal(t, tt.out, user.Email)
		assert.NotEqual(t, uint(0), user.ID)
	}
}

// 新規ユーザー作成 バリデーションエラー
func TestUserCreateUserValidateError(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  string
	}{
		{""},
		{"0"},
		{"-"},
		{"aaaa@"},
		{"@aaaa@example.com"},
		{"abcratyart.com"},
		{"111"},
		{"111.jp"},
	}
	for _, tt := range tests {
		user, err := models.CreateUser(models.DB, tt.in)
		// バリデーションが正常に動作している
		assert.Error(t, err)
		assert.Nil(t, user)
	}
}

// 新規登録　ユニークキー重複登録エラー
func TestUserCreateUserDuplicate(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  string
	}{
		{"test1@example.com"},
	}
	for _, tt := range tests {
		user, err := models.CreateUser(models.DB, tt.in)
		// ユニークキー制約でduplicateエラー発生する
		if err == nil {
			t.Error("can not be detected duplicate entry")
		}
		assert.Nil(t, user)
	}
}

// IDで検索
func TestUser_GetById(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  uint
		out string
	}{
		{1, "test1@example.com"},
	}
	for _, tt := range tests {
		user := models.User{}
		if err := user.GetById(models.DB, tt.in); err != nil{
			t.Error(err)
		}
		if user.Email != tt.out {
			assert.Equal(t, tt.out, user.Email)
		}
	}
	}

// IDで検索　対象レコード無し
func TestUser_GetByIdNotFound(t *testing.T)  {
	var assertT = assert.New(t)
	tests := []struct {
		in  uint
	}{
		{999},
		{12345},
	}
	for _, tt := range tests {
		user := models.User{}
		err := user.GetById(models.DB, tt.in);
		assertT.Equal(gorm.ErrRecordNotFound, err)
		assert.Equal(t, "", user.Email)
		assert.Equal(t, tt.in, user.ID)
	}
}

// IDで検索　エラー
func TestUser_GetByIdError(t *testing.T)  {
	var assertT = assert.New(t)
	tests := []struct {
		in  uint
	}{
		{0},
	}
	for _, tt := range tests {
		user := models.User{}
		err := user.GetById(models.DB, tt.in);
		assertT.Error(err)
		assertT.Equal("", user.Email)
		assertT.Equal(tt.in, user.ID)
	}
}

// IDで検索(排他ロック)
func TestUser_GetByIDForUpdate(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  uint
		out string
	}{
		{1, "test1@example.com"},
	}
	for _, tt := range tests {
		err := models.Transact(models.DB, func(db *gorm.DB) error {
			user := models.User{}
			if err := user.GetByIDForUpdate(models.DB, tt.in); err != nil{
				return err
			}
			if user.Email != tt.out {
				assert.Equal(t, tt.out, user.Email)
			}
			return nil
		})
		assert.Nil(t, err)
	}
}

// IDで検索　対象レコード無し(排他ロック)
func TestUser_GetByIDForUpdateNotFound(t *testing.T)  {
	var assertT = assert.New(t)
	tests := []struct {
		in  uint
	}{
		{999},
		{12345},
	}
	for _, tt := range tests {
		err := models.Transact(models.DB, func(db *gorm.DB) error {
			user := models.User{}
			err := user.GetByIDForUpdate(models.DB, tt.in);
			assertT.Equal(gorm.ErrRecordNotFound, err)
			assert.Equal(t, "", user.Email)
			assert.Equal(t, tt.in, user.ID)
			return err
		})
		assert.Error(t, err)
	}
}

// IDで検索　エラー(排他ロック)
func TestUser_GetByIDForUpdateError(t *testing.T)  {
	var assertT = assert.New(t)
	tests := []struct {
		in  uint
	}{
		{0},
	}
	for _, tt := range tests {
		err := models.Transact(models.DB, func(db *gorm.DB) error {
			user := models.User{}
			err := user.GetByIDForUpdate(models.DB, tt.in);
			assertT.Equal("", user.Email)
			assertT.Equal(tt.in, user.ID)
			return err
		})
		assert.Error(t, err)
	}
}

// Eメール検索
func TestUser_GetByEmail(t *testing.T) {
	tests := []struct {
		in  string
		out bool
	}{
		{"test1@example.com", true},
		{"earerefaewta2@example.com", false},
	}
	for _, tt := range tests {
		user := models.User{}
		err := user.GetByEmail(models.DB, tt.in)
		if tt.out {
			// 存在するメールアドレス
			assert.NoError(t, err)
			assert.NotEqual(t, uint(0), user.ID)
		} else {
			// 存在しないメールアドレス
			assert.Equal(t, gorm.ErrRecordNotFound, err)
			assert.Equal(t, uint(0), user.ID)
		}
	}
}

// ユーザー削除(チェックの整合性の為トランザクション化)
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
		err := models.Transact(models.DB, func(tx *gorm.DB) error {
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
func TestUser_DeleteUserByIdErrorTransaction(t *testing.T) {
	tests := []struct {
		in  uint
	}{
		{0}, // ID指定が不正
	}
	for _, tt := range tests {
		err := models.Transact(models.DB, func(tx *gorm.DB) error {
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

