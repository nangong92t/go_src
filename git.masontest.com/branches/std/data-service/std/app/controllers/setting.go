package controllers

import (
    // "std/data-service/std/app/modules"
	r "git.masontest.com/branches/goserver/workers/revel"
)

type Setting struct {
	AbstractController
}

// 设置通知条件.
func (c Setting) Notification(sid string, mine int, other int) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    return c.Render(true, nil)
}

