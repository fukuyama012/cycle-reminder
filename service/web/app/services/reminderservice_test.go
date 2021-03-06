package services_test

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// リマインド設定と紐付くリマインド予定をトランザクション作成
func TestCreateReminderSettingWithRelationInTransact(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
		BasisDate time.Time
	}{
		{1, "test name", "test title", "test text", 1, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
	}
	for _, tt := range tests {
		err := services.CreateReminderSettingWithRelationInTransact(models.DB, tt.UserID, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays, tt.BasisDate)
		assert.Nil(t, err)
	}
}

// リマインド設定と紐付くリマインド予定をトランザクション作成エラー
func TestCreateReminderSettingWithRelationInTransactError(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		UserID uint
		Name string
		NotifyTitle string
		NotifyText string
		CycleDays uint
		BasisDate time.Time
	}{
		// name無し
		{1, "", "test title", "test text", 1, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
		//　NotifyText無し
		{1, "test name2", "test title2", "", 365, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
		// User無し
		{9999, "test name2", "", "test text2", 7, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
		// CycleDaysが0
		{2, "test name2", "title", "test text2", 0, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
	}
	for _, tt := range tests {
		err := services.CreateReminderSettingWithRelationInTransact(models.DB, tt.UserID, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays, tt.BasisDate)
		assert.Error(t, err)
	}
}

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
		{2, "test name2", "title", "test text2", 31, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.February, 1, 0, 0, 0, 0, models.GetJSTLocation())},
		{1, "test name2", "", "test text2", 7, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()),
			time.Date(2018, time.January, 8, 0, 0, 0, 0, models.GetJSTLocation())},
	}
	err := models.Transact(models.DB, func(tx *gorm.DB) error {
		for _, tt := range tests {
			rSet, err := services.CreateReminderSettingWithRelation(tx, tt.UserID, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays, tt.BasisDate)
			if err != nil {
				return err
			}
			// リレーション情報としてリマインド予定にレコード追加されている
			rSch := models.ReminderSchedule{}
			errSch := rSch.GetByReminderSetting(tx, *rSet)
			assert.Nil(t, errSch)
			assert.NotEqual(t, uint(0), rSch.ID)
			// 次回通知日時が適切に設定されている
			assert.Equal(t, tt.NextNotifyDate, rSch.NotifyDate)

			// リマインダーが正常に設定されている
			assert.NotNil(t, rSet)
			assert.NotEqual(t, uint(0), rSet.ID)
			assert.Equal(t, tt.Name ,rSet.Name)
		}
		return nil
	})
	assert.Nil(t, err)
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
		data, err := models.TransactAndReceiveData(models.DB, func(tx *gorm.DB) (i interface{}, e error) {
			return services.CreateReminderSettingWithRelation(tx, tt.UserID, tt.Name, tt.NotifyTitle, tt.NotifyText, tt.CycleDays, tt.BasisDate)
		})
		assert.Error(t, err)
		// リマインダーが正常に設定されていない
		assert.Nil(t, data)
	}
}

// GetReminderListByUserID ユーザー情報からリマインド一覧取得
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
		reminderList, err := services.GetReminderListByUserID(models.DB, tt.UserID, tt.Limit, tt.Offset)
		assert.Nil(t, err)
		// limitとoffsetの兼ね合いで最大数決まる
		assert.Equal(t, tt.OutLen, len(reminderList))
		assert.Equal(t, tt.UserID, reminderList[0].UserID)
		assert.Equal(t, tt.Name, reminderList[0].Name)
		assert.Equal(t, tt.CycleDays, reminderList[0].CycleDays)
		assert.Equal(t, tt.NotifyDate, reminderList[0].NotifyDate)
		assert.NotZero(t, reminderList[0].ScheduleID)
	}
}

// GetReminderListByUserID ユーザー情報からリマインド一覧取得
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
		reminderList, errList := services.GetReminderListByUserID(models.DB, tt.UserID, tt.Limit, tt.Offset)
		assert.Nil(t, errList)
		// 存在しないUserで研削した場合は空
		assert.Equal(t, 0, len(reminderList))
	}
}

// GetReminderListByUserID ユーザー情報からリマインド一覧取得エラー
func TestGetReminderListByUserEmptyUser(t *testing.T) {
	prepareTestDB()
	user := models.User{}
	// 空User時はエラー
	reminderList, errList := services.GetReminderListByUserID(models.DB, user.ID, 1, 0)
	assert.Error(t, errList)
	assert.Nil(t, reminderList)
}

// GetReminderSchedulesReachedNotifyDate 通知日付に達した全リマインド予定の通知内容取得
func TestGetRemindersReachedNotifyDate(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		OutLen int
		SettingID uint
		ScheduleID uint
		Email string
		TargetDate time.Time
		OverID uint
		Limit int
	}{
		// limit, offset 別に正常系テスト
		{1, 1, 1, "test1@example.com", time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 10},
		{1, 1, 1, "test1@example.com", time.Date(2019, time.February, 27, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 10},
		{2, 1, 1, "test1@example.com",  time.Date(2019, time.February, 28, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 10},
		{2, 1, 1, "test1@example.com",  time.Date(2020, time.December, 30, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 10},
		{3, 1, 1, "test1@example.com",  time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 10},
		// 指定IDによる変化
		{2, 2, 2, "test1@example.com",  time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), uint(1), 10},
		{1, 3, 3, "test2@example.com",  time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), uint(1)+uint(1), 10},
		// limit変化
		{2, 1, 1, "test1@example.com",  time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 2},
		{1, 1, 1,  "test1@example.com",  time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 1},
	}
	for _, tt := range tests {
		reminderList, err := services.GetRemindersReachedNotifyDate(models.DB, tt.TargetDate, tt.OverID, tt.Limit)
		assert.Nil(t, err)
		// limitとoffsetの兼ね合いで最大数決まる
		assert.Equal(t, tt.OutLen, len(reminderList))
		assert.Equal(t, tt.Email, reminderList[0].Email)
		assert.Equal(t, tt.SettingID, reminderList[0].SettingID)
		assert.Equal(t, tt.ScheduleID, reminderList[0].ScheduleID)
	}
}

// GetReminderSchedulesReachedNotifyDate 通知日付に達した全リマインド予定の通知内容取得
// 該当無し
func TestGetRemindersReachedNotifyDateNotFound(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		OutLen int
		TargetDate time.Time
		OverID uint
		Limit int
	}{
		// ターゲット日付より古い通知日付のレコードなし
		{0, time.Date(2017, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), uint(0), 10},
		// ターゲット日付の指定IDを超えるレコードなし
		{0, time.Date(2020, time.December, 31, 0, 0, 0, 0, models.GetJSTLocation()), uint(3), 10},
	}
	for _, tt := range tests {
		reminderList, err := services.GetRemindersReachedNotifyDate(models.DB, tt.TargetDate, tt.OverID, tt.Limit)
		assert.Nil(t, err)
		assert.Equal(t, tt.OutLen, len(reminderList))
	}
}

// ResetReminderScheduleAfterNotify メール通知完了後の次回通知予定設定
func TestResetReminderScheduleAfterNotify(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		ReminderSettingID uint
		BasisDate time.Time
		NotifyDate string
	}{
		{1, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), "2018-01-08"},// 7日後
		{2, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation()), "2018-01-31"},// 30日後
	}
	for _, tt := range tests {
		err := models.Transact(models.DB, func(tx *gorm.DB) error {
			rSet := models.ReminderSetting{}
			if err := rSet.GetById(tx, tt.ReminderSettingID); err != nil {
				return err
			}
			err := services.ResetReminderScheduleAfterNotify(tx, rSet, tt.BasisDate)

			rSetCheck := models.ReminderSetting{}
			errSet := rSetCheck.GetById(tx, tt.ReminderSettingID)
			assert.Nil(t, errSet)
			// リマインド予定の次回通知日付が正しい
			rSch := models.ReminderSchedule{}
			errSch := rSch.GetByReminderSetting(tx, rSetCheck)
			assert.Nil(t, errSch)
			assert.Equal(t, tt.NotifyDate, rSch.NotifyDate.Format("2006-01-02"))
			return err
		})
		assert.Nil(t, err)
	}
}

// ResetReminderScheduleAfterNotify メール通知完了後の次回通知予定設定
// 対象レコード無し
func TestResetReminderScheduleAfterNotifyRecordNotFound(t *testing.T) {
	prepareTestDB()
	tests := []struct {
		ReminderSettingID uint
		BasisDate time.Time
	}{
		{99999, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
		{0, time.Date(2018, time.January, 1, 0, 0, 0, 0, models.GetJSTLocation())},
	}
	for _, tt := range tests {
		err := models.Transact(models.DB, func(tx *gorm.DB) error {
			rSet := models.ReminderSetting{}
			errSch := rSet.GetById(tx, tt.ReminderSettingID)
			assert.Error(t, errSch)

			err := services.ResetReminderScheduleAfterNotify(tx, rSet, tt.BasisDate)
			assert.Equal(t, gorm.ErrRecordNotFound, err)
			return err
		})
		assert.Error(t, err)
	}
}
