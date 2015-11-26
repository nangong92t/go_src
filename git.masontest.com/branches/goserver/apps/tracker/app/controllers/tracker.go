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
