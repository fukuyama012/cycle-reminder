package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
	"strconv"
	"time"
)

// Reminders is App
type Reminders struct {
	App
}

// Index リマインド一覧表示
func (c Reminders) Index() revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}

	rlist, err := services.GetReminderListByUserID(services.GetDB(), userID, 100, 0)
	if err != nil {
		c.Log.Errorf("Index GetReminderListByUserID %#v", err)
	}
	return c.Render(userID, rlist)
}

// UpdatePrepare リマインダー変更入力画面
func (c Reminders) UpdatePrepare(number int) revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}

	rSet, err := services.GetReminderSettingByUserIDAndNumber(services.GetDB(), userID, uint(number))
	if err != nil {
		// 変更対象存在しない
		c.Log.Errorf("UpdatePrepare() GetReminderSettingByUserIDAndNumber %#v", err)
		c.Redirect(routes.Reminders.Index())
	}
	return c.Render(userID, rSet)
}

// Update リマインダー変更
func (c Reminders) Update(number int) revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}

	// TODO 後ほど整理する。。
	result := "変更失敗！"
	cycleDays, errCast := strconv.Atoi(c.Params.Get("cycle_days"))
	if errCast != nil {
		c.Log.Errorf("cant cask cycleDays %#v", errCast)
		return c.Render(result)
	}
	// 変更処理
	_, err := services.UpdateReminderSettingByUserIDAndNumber(services.GetDB(), userID, uint(number), c.Params.Get("name"),
		c.Params.Get("notify_title"), c.Params.Get("notify_text"), uint(cycleDays))
	if err != nil {
		c.Log.Errorf("Update() UpdateReminderSettingByUserIDAndNumber %#v", err)
		return c.Render(result)
	}
	// リスト画面へ
	return c.Redirect(routes.Reminders.Index())
}

// CreatePrepare リマインド作成入力画面
func (c Reminders) CreatePrepare() revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	return c.Render(userID)
}

// Create リマインド作成
func (c Reminders) Create() revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}

	// TODO 後ほど整理する。。
	result := "登録失敗！"
	cycleDays, errCast := strconv.Atoi(c.Params.Get("cycle_days"))
	if errCast != nil {
		c.Log.Errorf("cant cask cycle_days %#v", errCast)
		return c.Render(result)
	}
	// 登録処理
	if err := services.CreateReminderSettingWithRelationInTransact(services.GetDB(), userID, c.Params.Get("name"),
		c.Params.Get("notify_title"), c.Params.Get("notify_text"), uint(cycleDays), time.Now());
	err != nil {
		c.Log.Errorf("CreateReminderSettingWithRelationInTransact %#v", err)
		return c.Render(result)
	}
	// （成功）リスト画面へ
	return c.Redirect(routes.Reminders.Index())
}

// Delete リマインダー削除
func (c Reminders) Delete(number int) revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}

	// TODO 後ほど整理する。。
	result := "削除失敗！"
	if err := services.DeleteReminderSettingByUserIDAndNumber(services.GetDB(), userID, uint(number)); err != nil {
		c.Log.Errorf("Delete() DeleteReminderSettingByUserIDAndNumber %#v", err)
		return c.Render(result)
	}
	// （成功）リスト画面へ
	return c.Redirect(routes.Reminders.Index())
}

