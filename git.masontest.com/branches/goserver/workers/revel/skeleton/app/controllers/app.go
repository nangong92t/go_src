package controllers

import "git.masontest.com/branches/goserver/workers/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}
