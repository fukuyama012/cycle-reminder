package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
)

type Login struct {
	*revel.Controller
}

// TOPページ　セッション有ればアプリへ 無ければLPへ
func (c Login) Index() revel.Result {
	return c.Render()
}

// oauth認証する
func (c Login) Oauth() revel.Result {
	state := services.RandString(10)
	
	session := services.Session{}
	session.SetForCallBackCheck(state)
	
	url := services.GetAuthCodeUrlWithSessionState(state)
	return c.Redirect(url)
}

// GoogleからのCallback処理 成功時セッション保存しアプリへ
func (c Login) Callback() revel.Result {
	return c.Render()
}

// ログアウト処理　セッション削除
func (c Login) Logout() revel.Result {
	return c.Render()
}



