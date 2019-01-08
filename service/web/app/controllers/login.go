package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
)

type Login struct {
	*revel.Controller
}

// TOPページ　セッション有ればアプリへ 無ければLPへ
func (c Login) Index()  {

}

// oauth認証する
func (c Login) Oauth() {
	state := services.RandString(10)
	
	session := services.Session{}
	session.SetForCallBackCheck(state)
	
	url := services.GetAuthCodeUrlWithSessionState(state)
	c.Redirect(url)
}

// GoogleからのCallback処理 成功時セッション保存しアプリへ
func (c Login) Callback()  {
	
}

// ログアウト処理　セッション削除
func (c Login) Logout()  {

}



