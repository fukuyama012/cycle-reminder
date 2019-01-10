package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
)

const googleOauthSession = "google_oauth_session"

type Auth struct {
	*revel.Controller
}

// TOPページ　セッション有ればアプリへ 無ければLPへ
func (c Auth) Index() revel.Result {
	return c.Render()
}

// oauth認証する
func (c Auth) Oauth() revel.Result {
	key := services.RandString(10)

	c.Session[googleOauthSession] = key
	url := services.GetAuthCodeUrlWithSessionKey(key)
	return c.Redirect(url)
}

// GoogleからのCallback処理 成功時セッション保存しアプリへ
func (c Auth) Callback() revel.Result {
	return c.Render()
}

// ログアウト処理　セッション削除
func (c Auth) Logout() revel.Result {
	return c.Render()
}



