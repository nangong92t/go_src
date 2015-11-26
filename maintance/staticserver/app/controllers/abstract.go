// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package controllers

import (
    "github.com/revel/revel"
    "maintance/staticserver/app/components/results"
)

// 抽象控制器, 所有具体实现控制器必须继承于此.
type AbstractController struct {
    *revel.Controller
}

// Overload the parent RenderError feature.
func (c *AbstractController) RenderError(errCode []interface{}) revel.Result {
    // return ErrorResult{c.RenderArgs, err}
    return c.RenderJson(results.RenderApiErrResult{Errors: results.GetError(errCode)})
}

// Overload the parent Render feature.
func (c *AbstractController) Render(data interface{}, err error) revel.Result {

    if err != nil {
        return c.RenderError([]interface{}{err})
    }

    return c.RenderJson(results.RenderApiResult{
        Result: data,
    })
}

func (c *AbstractController) NoFound() revel.Result {
    return c.RenderError([]interface{}{10012})
}

