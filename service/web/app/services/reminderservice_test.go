package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// リマインド設定と紐付くリマインド予定を作成
func TestCreateReminderSettingWithRelation(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
		BasisDate time.Time
		NextNotifyDate time.Time
	}{
		{1, "test name", "test title", "test text", 1, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.January, 2, 0, 0, 0, 0, models.GetJSTLocation())},
		{1, "test name2", "test title2", "test text2", 365, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2019, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
		{1, "test name2", "", "test text2", 7, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.January, 8, 0, 0, 0, 0, models.GetJSTLocation())},
		{2, "test name2", "title", "test text2", 31, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.February, 1, 0, 0, 0, 0, models.GetJSTLocation())},
	}
	for _, tt := range tests {
		rSet, err := services.CreateReminderSettingWithRelation(tt.UserID, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays, tt.BasisDate)
		assert.Nil(t, err)
		// リマインダーが正常に設定されている
		assert.NotNil(t, rSet)
		assert.NotEqual(t, uint(0), rSet.ID)
		assert.Equal(t, tt.Name ,rSet.Name)

		// リレーション情報としてリマインド予定にレコード追加されている
		rSch := models.ReminderSchedule{}
		errSch := rSch.GetByReminderSetting(models.DB, *rSet)
		assert.Nil(t, errSch)
		assert.NotEqual(t, uint(0), rSch.ID)
		// 次回通知日時が適切に設定されている
		assert.Equal(t, tt.NextNotifyDate, rSch.NotifyDate)
	}
}

// リマインド設定と紐付くリマインド予定を作成エラー
func TestCreateReminderSettingWithRelationError(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
		BasisDate time.Time
		NextNotifyDate time.Time
	}{
		// name無し
		{1, "", "test title", "test text", 1, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.January, 2, 0, 0, 0, 0, models.GetJSTLocation())},
			//　NotifyText無し
		{1, "test name2", "test title2", "", 365, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2019, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
			// User無し
		{9999, "test name2", "", "test text2", 7, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.January, 8, 0, 0, 0, 0, models.GetJSTLocation())},
			// CycleDaysが0
		{2, "test name2", "title", "test text2", 0, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.February, 1, 0, 0, 0, 0, models.GetJSTLocation())},
	}
	for _, tt := range tests {
		rSet, err := services.CreateReminderSettingWithRelation(tt.UserID, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays, tt.BasisDate)
		assert.Error(t, err)
		// リマインダーが正常に設定されていない
		assert.Nil(t, rSet)
	}
}

// GetReminderListByUser ユーザー情報からリマインド一覧取得
func TestGetReminderListByUser(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		Name string
		CycleDays uint
		NotifyDate time.Time
		Limit int
		Offset int
		OutLen int
	}{
		// limit, offset 別に正常系テスト
		{1, "name", 7, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), 1, 0, 1},
		{1, "name", 7, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), 2, 0, 2},
		{1, "name", 7, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), 3, 0, 3},
		{1, "name2", 30, time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), 3, 1, 2},
		{1, "name2", 30, time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), 2, 1, 2},
		{1, "name2", 30, time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), 1, 1, 1},
		{1, "name3", 60, time.Date(2099, time.July, 15, 0, 0, 0, 0, models.GetJSTLocation()), 3, 2, 1},
		{1, "name3", 60, time.Date(2099, time.July, 15, 0, 0, 0, 0, models.GetJSTLocation()), 2, 2, 1},
		{1, "name3", 60, time.Date(2099, time.July, 15, 0, 0, 0, 0, models.GetJSTLocation()), 1, 2, 1},
	}
	for _, tt := range tests {
		user := models.User{}
		err := user.GetById(models.DB, tt.UserID)
		assert.Nil(t, err)

		reminderList, errList := services.GetReminderListByUser(user, tt.Limit, tt.Offset)
		assert.Nil(t, errList)
		// limitとoffsetの兼ね合いで最大数決まる
		assert.Equal(t, tt.OutLen, len(reminderList))
		assert.Equal(t, tt.UserID, reminderList[0].UserID)
		assert.Equal(t, tt.Name, reminderList[0].Name)
		assert.Equal(t, tt.CycleDays, reminderList[0].CycleDays)
		assert.Equal(t, tt.NotifyDate, reminderList[0].NotifyDate)
	}
}

// GetReminderListByUser ユーザー情報からリマインド一覧取得
func TestGetReminderListByUserNotExistsUser(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		Limit int
		Offset int
	}{
		{99999, 1, 0},
	}
	for _, tt := range tests {
		user := models.User{}
		err := user.GetById(models.DB, tt.UserID)
		assert.Error(t, err)

		reminderList, errList := services.GetReminderListByUser(user, tt.Limit, tt.Offset)
		assert.Nil(t, errList)
		// 存在しないUserで研削した場合は空
		assert.Equal(t, 0, len(reminderList))
	}
}

// GetReminderListByUser ユーザー情報からリマインド一覧取得エラー
func TestGetReminderListByUserEmptyUser(t *testing.T) {
	prepareTestDB()
	user := models.User{}
	// 空User時はエラー
	reminderList, errList := services.GetReminderListByUser(user, 1, 0)
	assert.Error(t, errList)
	assert.Nil(t, reminderList)
}