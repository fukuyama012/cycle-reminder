package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
)

const (
	googleOauthSession  = "google_oauth_session"
	serviceLoginSession = "service_login_session"
)

type Auth struct {
	*revel.Controller
}

// TOPページ　セッション有ればアプリへ 無ければLPへ
func (c Auth) Index() revel.Result {
	// loginチェック
	_, err := c.Session.Get(serviceLoginSession);
	if err != nil {
		// セッション無ければLPへ
		return c.Redirect(routes.App.Index())
	}
	return c.Redirect(routes.Reminders.Index())
}

// oauth認証する
func (c Auth) Login() revel.Result {
	key := services.RandString(10)

	c.Session[googleOauthSession] = key
	url := services.GetAuthCodeUrlWithSessionKey(key)
	return c.Redirect(url)
}

// GoogleからのCallback処理 成功時セッション保存しアプリへ
func (c Auth) Callback() revel.Result {
	if !c.isValidCallbackSession() {
		return c.Redirect(routes.App.Index())
	}

	code := c.Params.Query.Get("code");
	if code == "" {
		c.Log.Error("not exists Callback Code")
		return c.Redirect(routes.App.Index())
	}

	oauthInfo, err := services.GetOauthInfo(code)
	if err != nil {
		c.Log.Errorf("oauth info, %#v", err)
		return c.Redirect(routes.App.Index())
	}

	userId, err := services.GetUserIdOrCreateUserId(oauthInfo.Email)
	if err != nil {
		c.Log.Errorf("get or create userId, %#v", err)
		return c.Redirect(routes.App.Index())
	}
	
	c.Session[serviceLoginSession] = userId
	return c.Redirect(routes.Reminders.Index())
}

// セッションを通じて正当なリクエストかチェック
func (c Auth) isValidCallbackSession() bool  {
	state := c.Params.Query.Get("state");
	if state == "" {
		c.Log.Error("not exists Callback Session")
		return false
	}
	oauthSession, ok := c.Session[googleOauthSession];
	if !ok {
		c.Log.Error("not exists Oauth Session")
		return false
	}
	if oauthSession.(string) != state {
		c.Log.Error("invalid Callback Session")
		return false
	}
	return true
}

// ログアウト処理　セッション削除
func (c Auth) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}



