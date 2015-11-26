package controllers

import (
    r "git.masontest.com/branches/goserver/workers/revel"
    "std/data-service/std/app/modules"
    "std/data-service/std/app/models"
    "errors"
    "time"
)

type Label struct {
    AbstractController
}

// 获取label列表.
func (c Label) GetList(page int, limit int) r.Result {
    labelModule := modules.NewLabel(c.Db)
    list, err   := labelModule.GetList(page, limit, false)

    return c.Render(list, err)
}

// 设置用户喜欢的标签.
func (c Label) Like(sid string, ids string) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    labelModule := modules.NewLabel(c.Db)
    isOk        := labelModule.Like(session.SessionVal.UserId, ids)

    return c.Render(isOk, nil)
}

func (c Label) Add(sid, label string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    if label == "" {
        return c.Render(nil, errors.New("sorry, no label!"))
    }

    isExists    := modules.NewLabel(c.Db).GetByName(label)
    if isExists.Id != 0 {
        return c.Render(nil, errors.New("Sorry, this label is exists!"))
    }

    labelD   := &models.Label {
        Name: label,
        Created: time.Now().Unix(),
    }

    err     := c.Db.Insert(labelD)

    if err != nil {
        return c.Render(nil,err)
    }

    return c.Render(labelD, nil)
}
