package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
	"os"
)

// App Revel Controller
type App struct {
	*revel.Controller
	UserInfo struct {
		ID uint
		IsLogin bool
	}
}

// Index TOPページ
func (c App) Index() revel.Result {
	return c.Render()
}

// Terms 利用規約
func (c App) Terms() revel.Result {
	return c.Render()
}

// getLoginInfo ログイン情報取得
func (c *App) getLoginInfo() revel.Result {
	c.UserInfo.IsLogin = false
	c.UserInfo.ID = uint(0)
	// loginチェック
	userIDSession, ok := c.getUserIDBySession()
	if !ok {
		// 有効なセッション情報無し
		return nil
	}
	user, err := services.CheckUserID(userIDSession)
	if err != nil {
		// ユーザー登録無し
		return nil
	}
	c.UserInfo.IsLogin = true
	c.UserInfo.ID = user.ID
	return nil
}

// getUserIDBySession セッションからユーザーID取り出し
func (c App) getUserIDBySession() (uint, bool) {
	userIDSesson, err := c.Session.Get(serviceLoginSession);
	if err != nil {
		return 0, false
	}
	val, ok := userIDSesson.(float64)
	if !ok {
		return 0, false
	}
	return uint(val), true
}

// setCommonToView　Viewに共通情報を設定
func (c App) setCommonToView() revel.Result {
	c.ViewArgs["googleAnalytics"] = os.Getenv("GOOGLE_ANALYTICS")
	c.ViewArgs["isLogin"] = c.UserInfo.IsLogin
	return nil
}