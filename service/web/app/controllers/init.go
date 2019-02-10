package controllers

import "github.com/revel/revel"

func init()  {
	revel.InterceptMethod((*App).getLoginInfo, revel.BEFORE)
	revel.InterceptMethod(App.setCommonToView, revel.BEFORE)
}
