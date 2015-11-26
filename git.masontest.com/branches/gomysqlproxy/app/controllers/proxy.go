// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package controllers

import (
    "git.masontest.com/branches/goserver/workers/revel"
    "git.masontest.com/branches/gomysqlproxy/app/models"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
)

type Proxy struct {
    *revel.Controller
}

func (c *Proxy) Exec(args ...string) revel.Result {
    user    := client.NewClient("58.96.178.236")

    res, err    := models.MyProxy.Exec(args, user)
    return c.Render(res, err)
}

func (c *Proxy) Getstatus() revel.Result {
    return c.Render(models.MyProxy.GetStatus())
}
