// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package controllers

import (
    "git.masontest.com/branches/goserver/workers/revel"
)

// The Host checker.
type Tracker struct {
    *revel.Controller
}

func (c *Tracker) Status() revel.Result {
    return c.Render("OK", nil)
}
