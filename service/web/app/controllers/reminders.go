package controllers

import (
	"github.com/fukuyama012/cycle-reminder/service/web/app/routes"
	"github.com/fukuyama012/cycle-reminder/service/web/app/services"
	"github.com/revel/revel"
	"strconv"
	"time"
)

type Reminders struct {
	*revel.Controller
}

// Index リマインド一覧表示
func (c Reminders) Index() revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	
	rlist, err := services.GetReminderListByUserID(services.GetDB(), userID, 100, 0)
	if err != nil {
		c.Log.Errorf("Index GetReminderListByUserID %#v", err)
	}
	isLogin := true
	return c.Render(userID, isLogin, rlist)
}

// CreatePrepare リマインド作成入力画面
func (c Reminders) CreatePrepare() revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}
	isLogin := true
	return c.Render(userID, isLogin)
}

// Create リマインド作成
func (c Reminders) Create() revel.Result {
	// loginチェック
	userID := c.getLoginUser()
	if userID == uint(0) {
		// 未ログイン TOP LPへ
		return c.Redirect(routes.App.Index())
	}

	// TODO 後ほど整理する。。
	isLogin := true
	result := "登録失敗！"
	cycle_days, errCast := strconv.Atoi(c.Params.Get("cycle_days"))
	if errCast != nil {
		c.Log.Errorf("cant cask cycle_days %#v", errCast)
		return c.Render(result, isLogin)
	}
	// 登録処理
	_, err := services.CreateReminderSettingWithRelation(services.GetDB(), userID, c.Params.Get("name"),
		c.Params.Get("notify_title"), c.Params.Get("notify_text"), uint(cycle_days), time.Now())
	if err != nil {
		c.Log.Errorf("error, CreateReminderSettingWithRelation %#v", err)
		return c.Render(result, isLogin)
	}
	// リスト画面へ
	return c.Redirect(routes.Reminders.Index())
}

// getLoginUser ログインユーザー情報取得
func (c Reminders) getLoginUser() (uint) {
	// loginチェック
	userIdSession, ok := c.getUserIdBySession()
	if !ok {
		// 有効なセッション情報無し
		return uint(0)
	}
	user, err := services.CheckUserID(userIdSession)
	if err != nil {
		// ユーザー登録無し
		return uint(0)
	}
	return user.ID
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
