package controllers

import (
    r "github.com/revel/revel"
    "std/data-service/std/app/modules"
)

type Label struct {
    AbstractController
}

// 获取label列表.
func (c Label) GetList(page int, limit int) r.Result {
    labelModule := modules.NewLable(c.Db)
    list, err   := labelModule.GetList(page, limit)

    return c.Render(list, err)
}

// 设置用户喜欢的标签.
func (c Label) Like(sid string, ids string) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    labelModule := modules.NewLable(c.Db)
    isOk        := labelModule.Like(session.SessionVal.UserId, ids)

    return c.Render(isOk, nil)
}
