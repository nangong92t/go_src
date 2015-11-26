package controllers

import (
    "std/data-service/std/app/modules"
	r "git.masontest.com/branches/goserver/workers/revel"
)

type User struct {
	AbstractController
}

func (c User) GetUsers(userId int) r.Result {
	greeting := "It works!"

    userModule  := modules.NewUser(c.Db)
    users       := userModule.GetUsers()

    c.Session["user_name"]  =   greeting

	return c.RenderJson(users)
}

func (c User) GetUserByName(name string) r.Result {
    userModule  := modules.NewUser(c.Db)
    user        := userModule.GetUserByName(name)

	return c.Render(user, nil)
}

// 设置用户Profile.
func (c User) Setting(sid string, gender int, age int, lllness int) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    res, err    := modules.NewUser(c.Db).SettingProfile(session.SessionVal.UserId, gender, age, lllness)

    return c.Render(res, err)
}

// 获取当前profile所有者创建的topic列表. 
func (c User) GetTopicList(uid int64, page int, limit int) r.Result {
    topicModule := modules.NewTopic(c.Db, c.Rc)

    var topics *modules.Pager
    var err error

    topics, err = topicModule.GetTopicList(page, "", limit, 0, false)

    return c.Render(topics, err)
}

// 获取当前profile用户订阅Topic列表.
func (c User) GetSubscribeTopicList(uid int64, page int, limit int) r.Result {
    topicModule := modules.NewTopic(c.Db, c.Rc)

    var topics *modules.Pager
    var err error

    topics, err = topicModule.GetTopicList(page, "", limit, 0, false)

    return c.Render(topics, err)
}

// 拉黑某用户
func (c User) Block(sid string, blocked int64) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    err := modules.NewUser(c.Db).Block(session.SessionVal.UserId, blocked, c.Sess.IsAdmin())

    return c.Render(true, err)
}

// 获取我的黑名单列表
func (c User) GetBlockedUsers(sid string, page int, limit int) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    blockedUsers, err   := modules.NewUser(c.Db).GetBlockedUsers(session.SessionVal.UserId, page, limit)

    return c.Render(blockedUsers, err)
}
