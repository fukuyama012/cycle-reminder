package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/revel/revel"
)

type Reminders struct {
	Auth
}

func (c Reminders) Index() revel.Result {
	// loginチェック
	userId, err := c.Session.Get(serviceLoginSession);
	if err != nil {
		// セッション無ければLPへ
		return c.Redirect(routes.App.Index())
	}
	// TODO 登録情報チェック
	return c.Render(userId)
}
