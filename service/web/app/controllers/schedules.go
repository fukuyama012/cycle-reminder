package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/revel/revel"
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
	// TODO idからsettingsをリレーションから取得してuserIDが一致するかチェック

	return c.Render()
}

// Update リマインダー変更
func (c Schedules) Update(number int) revel.Result {
	if !c.UserInfo.IsLogin {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	// TODO idからsettingsをリレーションから取得してuserIDが一致するかチェック
	// TODO ファイル内共通処理化できる？

	// リスト画面へ
	return c.Redirect(routes.Reminders.Index())
}
