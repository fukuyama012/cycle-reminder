package models_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/jinzhu/gorm"
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
		rSet, err := models.CreateReminderSetting(models.DB, user, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
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
		rSet, err := models.CreateReminderSetting(models.DB, tt.user, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		// リマインダーが正常に設定されている
		assert.Error(t, err)
		assert.Nil(t, rSet)
	}
}

// Userで検索
func TestGetReminderSettingsByUser(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		UserIDOut uint
	}{
		{1, 1},
		{2, 2},
	}
	for _, tt := range tests {
		user := models.User{}
		if err := user.GetById(tt.UserID); err != nil {
			t.Error(err)
		}
		rSettings, err := models.GetReminderSettingsByUser(models.DB, user)
		if err != nil {
			t.Error(err)
		}
		assert.NotNil(t, rSettings)
		for _, rs := range rSettings {
			assert.Equal(t, tt.UserIDOut, rs.UserID)
		}
	}
}

// ユーザー検索　対象レコード無し
func TestGetReminderSettingsByUserRecordNotFound(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
	}{
		{3},
	}
	for _, tt := range tests {
		user := models.User{}
		if err := user.GetById(tt.UserID); err != nil {
			t.Error(err)
		}
		rSettings, err := models.GetReminderSettingsByUser(models.DB, user)
		// gorm.ErrRecordNotFoundは返却されない
		assert.NoError(t, err)
		// モデルの空配列が返却される
		assert.Equal(t, 0, len(rSettings))
	}
}

// IDで検索
func TestReminderSetting_GetById(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		In  uint
		Name string
		CycleDays uint
	}{
		{1, "name", 7},
		{2, "name2", 100},
		{3, "name3", 1},
	}
	for _, tt := range tests {
		rs := models.ReminderSetting{}
		if err := rs.GetById(models.DB, tt.In); err != nil{
			t.Error(err)
		}
		assert.Equal(t, tt.Name, rs.Name)
		assert.Equal(t, tt.CycleDays, rs.CycleDays)
	}
}

// IDで検索　対象レコード無し
func TestReminderSetting_GetByIdRecordNotFound(t *testing.T) {
	tests := []struct {
		in  uint
	}{
		{999},
		{12345},
	}
	for _, tt := range tests {
		rs := models.ReminderSetting{}
		err := rs.GetById(models.DB, tt.in);
		assert.Equal(t, gorm.ErrRecordNotFound, err)
		assert.Equal(t, "", rs.NotifyTitle)
		assert.Equal(t, tt.in, rs.ID)
	}
}

// フィールド更新
func TestReminderSetting_Updates(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		ID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
	}{
		{1, "test name1", "test title", "test text", 1},
		{2, "テストネーム2", "テストタイトル2", "テストテキスト2", 365},
		{1, "test name3", "", "test text3", 10},
	}
	for _, tt := range tests {
		rSet := models.ReminderSetting{}
		if err := rSet.GetById(models.DB, tt.ID); err != nil {
			t.Error(err)
		}
		err := rSet.Updates(models.DB, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		assert.Nil(t, err)
		assert.Equal(t, tt.Name, rSet.Name)
		assert.Equal(t, tt.NotifyTitle, rSet.NotifyTitle)
		assert.Equal(t, tt.NotifyText, rSet.NotifyText)
		assert.Equal(t, tt.CycleDays, rSet.CycleDays)
	}
}

// フィールド更新 バリデーションエラー
func TestReminderSetting_UpdatesValidationError(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		ID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
	}{
		{1, "", "test title", "test text", 1}, // name無し
		// name 最大長超え
		{1, "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", "title", "test text2", 7},
		{1, "test name", "test title", "", 365}, // テキスト無し
		// タイトル最大長超え
		{1, "name", "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890", "test text2", 7},
		{1, "test name", "test title", "text", 0}, // リマインド日数0
		{1, "test name", "test title", "text", 366}, // リマインド日数最大値超え
	}
	for _, tt := range tests {
		rSet := models.ReminderSetting{}
		if err := rSet.GetById(models.DB, tt.ID); err != nil {
			t.Error(err)
		}
		err := rSet.Updates(models.DB, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		assert.Error(t, err)
	}
}

// フィールド更新 バリデーションエラー
// ID無し
func TestReminderSetting_UpdatesNoIdError(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		ID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
	}{
		{999, "test name", "test title", "test text", 1}, // ID無し
	}
	for _, tt := range tests {
		rSet := models.ReminderSetting{}
		if err := rSet.GetById(models.DB, tt.ID); err == nil {
			t.Error(err)
		}
		err := rSet.Updates(models.DB, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays)
		assert.Error(t, err)
	}
}

// 削除
func TestReminderSetting_DeleteById(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  uint
		out bool
	}{
		{1, true},
		{2, true},
		{9999, false},
	}
	for _, tt := range tests {
		recordCountBefore, errCount := models.CountReminderSetting(models.DB)
		if errCount != nil {
			t.Errorf("reminder setting count err %#v", errCount)
		}

		rs := models.ReminderSetting{}
		err := rs.DeleteById(models.DB, tt.in);
		if err != nil {
			t.Errorf("reminder setting Delete err %#v", err)
		}
		recordCountAfter, errCount := models.CountReminderSetting(models.DB)
		if errCount != nil {
			t.Errorf("reminder setting count err %#v", errCount)
		}
		if tt.out {
			// レコードが減少している
			assert.Equal(t, recordCountBefore - 1, recordCountAfter)
		} else {
			// 存在しないID
			// レコードが減少していない
			assert.Equal(t, recordCountBefore, recordCountAfter)
		}
	}
}

// 削除 id=0の場合は個別エラー（適切にエラー処理しないと全て削除される）
func TestReminderSetting_DeleteByIdZeroValueError(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		in  uint
	}{
		{0},
	}
	for _, tt := range tests {
		recordCountBefore, errCount := models.CountReminderSetting(models.DB)
		if errCount != nil {
			t.Errorf("reminder setting count err %#v", errCount)
		}

		rs := models.ReminderSetting{}
		err := rs.DeleteById(models.DB, tt.in);
		assert.Error(t, err)
		recordCountAfter, errCount := models.CountReminderSetting(models.DB)
		if errCount != nil {
			t.Errorf("reminder setting count err %#v", errCount)
		}
		// id=0指定エラー時
		// レコードが減少していない
		assert.Equal(t, recordCountBefore, recordCountAfter)
	}
}