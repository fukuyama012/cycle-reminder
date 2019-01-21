package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
)

type Reminders struct {
	Auth
}

func (c Reminders) Index() revel.Result {
	// loginチェック
	userIdSession, ok := c.getUserIdBySession()
	if !ok {
		// 有効なセッション情報無ければLPへ
		return c.Redirect(routes.App.Index())
	}
	user, err := services.CheckUserID(userIdSession)
	if err != nil {
		// ユーザー登録無ければLPへ
		return c.Redirect(routes.App.Index())
	}
	userId := user.ID
	return c.Render(userId)
}

func (c Reminders) getUserIdBySession() (uint, bool) {
	userIdSesson, err := c.Session.Get(serviceLoginSession);
	if err != nil {
		return 0, false
	}
	val, ok := userIdSesson.(float64)
	if !ok {
		return 0, false
	}
	return uint(val), true
}
