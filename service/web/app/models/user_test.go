package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		if err := user.GetById(tt.in); err != nil{
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
		err := user.GetById(tt.in); 
		assertT.Equal(gorm.ErrRecordNotFound, err)
		assert.Equal(t, "", user.Email)
		assert.Equal(t, tt.in, user.ID)
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
		err := user.GetByEmail(tt.in)
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

// 新規ユーザー作成
func TestUser_CreateUser(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"insert1@example.com", "insert1@example.com"},
		{"insert2@example.com", "insert2@example.com"},
	}
	for _, tt := range tests {
		user, err := models.CreateUser(tt.in)
		if err != nil {
			t.Error(err)
		}
		// メールアドレスが正常に設定されている
		assert.Equal(t, tt.out, user.Email)
		assert.NotEqual(t, uint(0), user.ID)
	}
}

// 新規登録　ユニークキー重複登録エラー
func TestUser_CreateUserDuplicate(t *testing.T) {
	tests := []struct {
		in  string
	}{
		{"test1@example.com"},
	}
	for _, tt := range tests {
		user, err := models.CreateUser(tt.in)
		// ユニークキー制約でduplicateエラー発生する
		if err == nil {
			t.Error("can not be detected duplicate entry")
		} 
		assert.Nil(t, user)
	}
}

/*
func TestUser_GetAll(t *testing.T) {
	prepareTestDB()
	users, err := models.GetAllUsers()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, len(users))
}
*/