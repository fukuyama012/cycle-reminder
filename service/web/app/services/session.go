package services

import "github.com/revel/revel"

const GOOGLE_OAUTH_STATE = "googleOauthState"

type Session struct {
	*revel.Controller
}

func (c Session) SetForCallBackCheck(state string) {
	c.Session[GOOGLE_OAUTH_STATE] = state
	c.Session.SetNoExpiration()
}
