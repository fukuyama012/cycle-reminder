package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
	"time"
)

// Reminders is App
type Schedules struct {
	App
}

// Index 画面無し　リマインダ一覧へ戻る
func (c Schedules) Index() revel.Result {
	if !c.UserInfo.IsLogin {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}

	return c.Redirect(routes.Auth.Index())
}

// UpdatePrepare リマインダー変更入力画面
func (c Schedules) UpdatePrepare(id int) revel.Result {
	if !c.UserInfo.IsLogin {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	
	rSet, rSch, err := services.GetSettingAndScheduleByScheduleIDAndUserID(services.GetDB(), uint(id), c.UserInfo.ID)
	if err != nil {
		c.Log.Errorf("not found schedule %#v", err)
		return c.Redirect(routes.Reminders.Index())
	}
	return c.Render(rSet, rSch)
}

// Update リマインダー変更
func (c Schedules) Update(id int) revel.Result {
	if !c.UserInfo.IsLogin {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	_, rSch, err := services.GetSettingAndScheduleByScheduleIDAndUserID(services.GetDB(), uint(id), c.UserInfo.ID)
	if err != nil {
		c.Log.Errorf("not found schedule %#v", err)
		return c.Redirect(routes.Reminders.Index())
	}

	jst, _ := time.LoadLocation("Asia/Tokyo")
	notify, err := time.ParseInLocation("2006-01-02 15:04:05", c.Params.Get("notify_date") + " 00:00:00", jst)
	if err != nil {
		c.Log.Errorf("cant cast notify date %#v", err)
		result := "正しい日付形式で入力してください"
		return c.Render(result)
	}
	now := time.Now()
	if notify.Before(now) {
		result := "本日より後の日付を入力してください"
		return c.Render(result)
	}
	if err := rSch.Updates(services.GetDB(), notify); err != nil {
		c.Log.Errorf("Update err %#v", err)
		result := "変更失敗！"
		return c.Render(result)
	}

	return c.Redirect(routes.Reminders.Index())
}
