package controllers

import (
//	rpc "git.masontest.com/branches/goserver/client"
    "git.masontest.com/branches/goserver/apps/hoststatus/app/models"
    "git.masontest.com/branches/goserver/workers/revel"
)

// The Host checker.
type Host struct {
    *revel.Controller
}

func (c *Host) Status() revel.Result {
    result, err := models.GetOSStatus()
    return c.Render(result, err)
}

func (c *Host) Ping() revel.Result {
    return c.Render("Pong", nil)
}

func (c *Host) PingTo(ip string) revel.Result {
    return c.Render(models.Ping(ip, 3))
}

func (c *Host) ToRun(cmd string, arg ...string) revel.Result {
    return c.Render(models.ToRun(cmd, arg...), nil)
}
