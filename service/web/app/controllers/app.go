package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
	"os"
)

// App Revel Controller
type App struct {
	*revel.Controller
}

// Index TOPページ
func (c App) Index() revel.Result {
	userID := c.getLoginUser()
	isLogin := false
	if userID != uint(0) {
		isLogin = true
	}
	return c.Render(isLogin)
}

// Terms 利用規約
func (c App) Terms() revel.Result {
	userID := c.getLoginUser()
	isLogin := false
	if userID != uint(0) {
		isLogin = true
	}
	return c.Render(isLogin)
}

// getLoginUser ログインユーザー情報取得
func (c App) getLoginUser() (uint) {
	// loginチェック
	userIDSession, ok := c.getUserIDBySession()
	if !ok {
		// 有効なセッション情報無し
		return uint(0)
	}
	user, err := services.CheckUserID(userIDSession)
	if err != nil {
		// ユーザー登録無し
		return uint(0)
	}
	return user.ID
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
	return nil
}