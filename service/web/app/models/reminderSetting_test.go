package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 新規ユーザー作成
func TestCreateReminderSetting(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
	}{
		{1, "test name", "test title", "test text", 1},
		{1, "test name2", "test title2", "test text2", 365},
		{1, "test name2", "", "test text2", 7},
		{2, "test name2", "title", "test text2", 7},
	}
	for _, tt := range tests {
		user := models.User{}
		if err :=user.GetById(tt.UserID); err != nil {
			t.Error(err)
		}
		rSet, err := models.CreateReminderSetting(user, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		if err != nil {
			t.Error(err)
		}
		// リマインダーが正常に設定されている
		assert.NotNil(t, rSet)
		assert.Equal(t, rSet.UserID, user.ID)
	}
}

// 新規ユーザー作成 エラー
func TestCreateReminderSettingError(t *testing.T) {
	user1 := models.User{}
	if err :=user1.GetById(1); err != nil {
		t.Error(err)
	}
	user2 := models.User{}
	tests := []struct {
		user models.User
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
	}{
		{user2, "name", "test title", "test text", 1}, // 空user
		{user1, "", "test title", "test text", 1}, // name無し
		// name 最大長超え
		{user1, "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", "title", "test text2", 7},
		{user1, "test name", "test title", "", 365}, // テキスト無し
		// タイトル最大長超え
		{user1, "name", "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", "test text2", 7},
		{user1, "test name", "test title", "text", 0}, // リマインド日数0
		{user1, "test name", "test title", "text", 366}, // リマインド日数最大値超え
	}
	for _, tt := range tests {
		rSet, err := models.CreateReminderSetting(tt.user, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		// リマインダーが正常に設定されている
		assert.Error(t, err)
		assert.Nil(t, rSet)
	}
} 