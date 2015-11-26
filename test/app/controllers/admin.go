package controllers

import (
    "github.com/revel/revel"
)

type Admin struct {
    GorpController
}

func (c Admin) GetUsers() revel.Result {
    return c.RenderText("Hello Tony")
}
