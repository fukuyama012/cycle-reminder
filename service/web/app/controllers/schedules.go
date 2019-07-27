package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/models"
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
	"time"
)

// Reminders is App
type Schedules struct {
	App
}

// UpdatePrepare リマインダー変更入力画面
func (c Schedules) UpdatePrepare(id int) revel.Result {
	if !c.UserInfo.IsLogin {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	rSch := models.ReminderSchedule{}
	if err := rSch.GetByID(services.GetDB(), uint(id)); err != nil {
		c.Log.Errorf("not found schedule %#v", err)
		c.Redirect(routes.Reminders.Index())
	}
	rSet := models.ReminderSetting{}
	if err := rSet.GetByIDAndUserID(services.GetDB(), rSch.ReminderSettingID, c.UserInfo.ID); err != nil {
		c.Log.Errorf("not found setting %#v", err)
		c.Redirect(routes.Reminders.Index())
	}
	return c.Render(rSet, rSch)
}

// Update リマインダー変更
func (c Schedules) Update(id int) revel.Result {
	if !c.UserInfo.IsLogin {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	rSch := models.ReminderSchedule{}
	if err := rSch.GetByID(services.GetDB(), uint(id)); err != nil {
		c.Log.Errorf("not found schedule %#v", err)
		c.Redirect(routes.Reminders.Index())
	}
	rSet := models.ReminderSetting{}
	if err := rSet.GetByIDAndUserID(services.GetDB(), rSch.ReminderSettingID, c.UserInfo.ID); err != nil {
		c.Log.Errorf("not found setting %#v", err)
		c.Redirect(routes.Reminders.Index())
	}

	result := "変更失敗！"
	date := c.Params.Get("notify_date") + " 00:00:00"
	notify, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		c.Log.Errorf("cant cast notify date %#v", err)
		return c.Render(result)
	}
	
	if err := rSch.Updates(services.GetDB(), notify); err != nil {
		c.Log.Errorf("Update err %#v", err)
		return c.Render(result)
	}

	return c.Redirect(routes.Reminders.Index())
}
