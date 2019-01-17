package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		if err := user.GetById(models.DB, tt.UserID); err != nil {
			t.Error(err)
		}
		rSet, err := services.CreateReminderSetting(user, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		assert.Nil(t, err)
		// リマインダーが正常に設定されている
		assert.NotNil(t, rSet)
		assert.Equal(t, rSet.UserID, user.ID)
	}
}

// 新規ユーザー作成 エラー
func TestCreateReminderSettingError(t *testing.T) {
	prepareTestDB()
	user1 := models.User{}
	if err :=user1.GetById(models.DB, 1); err != nil {
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
		{user1, services.RandString(101), "title", "test text2", 7}, // name 最大長超え
		{user1, "test name", "test title", "", 365}, // テキスト無し
		{user1, "name", services.RandString(101), "test text2", 7}, // タイトル最大長超え
		{user1, "test name", "test title", "text", 0}, // リマインド日数0
		{user1, "test name", "test title", "text", 366}, // リマインド日数最大値超え
	}
	for _, tt := range tests {
		rSet, err := services.CreateReminderSetting(tt.user, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		// リマインダーが正常に設定されていない
		assert.Error(t, err)
		assert.Nil(t, rSet)
	}
}
