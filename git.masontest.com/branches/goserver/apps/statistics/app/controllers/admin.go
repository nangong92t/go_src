package controllers

import (
    "std/data-service/std/app/modules"
    "std/data-service/std/app/helpers"
	r "github.com/revel/revel"
    "errors"
    "time"
)

type Admin struct {
	AbstractController
}

func (c Admin) GetUsers(userId int) r.Result {
	greeting := "It works!"

    userModule  := modules.NewUser(c.Db)
    users       := userModule.GetUsers()

    c.Session["user_name"]  =   greeting

	return c.Render(users, nil)
}

func (c Admin) GetAllStat(sid string, date string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    if date == "" {
        return c.Render(nil, errors.New("Sorry, no date."))
    }

    dateInt, err    := time.Parse("2006/01/02", date)
    if err != nil { return c.Render(nil, err) }
    dateUInt        := dateInt.Unix()

    userModule      := modules.NewUser(c.Db)
    topicModule     := modules.NewTopic(c.Db, c.Rc)

    // get user stat data.
    userStat, err   := userModule.GetStat(dateUInt)
    if err != nil { return c.Render(nil, err) }

    // get topic stat data.
    topicStat, err  := topicModule.GetStat(dateUInt)
    if err != nil { return c.Render(nil, err) }

    // get comment stat data.
    commentStat, err := topicModule.GetCommentStat(dateUInt)
    if err != nil { return c.Render(nil, err) }

    // get user date news count.
    userDateNews, err   := userModule.GetOneMonthDateNewCnt(dateUInt)
    if err != nil { return c.Render(nil, err) }

    // get topic date news count.
    topicDateNews, err  := topicModule.GetOneMonthDateNewCnt(dateUInt)
    if err != nil { return c.Render(nil, err) }

    // get comment date news count.
    commentDateNews, err := topicModule.GetCommOneMonthDateNewCnt(dateUInt)
    if err != nil { return c.Render(nil, err) }

    // response format
    response    := map[string]interface{}{
        "detail": map[string]interface{}{
            "users": userStat,
            "topics": topicStat,
            "comments": commentStat,
        },
        "date": map[string]interface{}{
            "users": userDateNews,
            "topics": topicDateNews,
            "comments": commentDateNews,
        },
    }

    return c.Render(response, nil)
}

func (c Admin) GetUserList(sid string, page, limit int, utype string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    userModule  := modules.NewUser(c.Db)
    var list modules.Pager
    var err error
    switch (utype) {
        case "active":
            list, err   = userModule.GetMostPostUserList(page, limit, 20)
        case "blocked":
            list, err   = userModule.GetBlockUserList(page, limit, true)
        case "block":
            list, err   = userModule.GetBlockUserList(page, limit, false)
        default:
            return c.Render(nil, errors.New("Sorry, the utype error!"))
    }

    return c.Render(list, err)
}

func (c Admin) RemoveUsers(sid string, ids string, withPost bool) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }
    uid := helpers.GetArrayFromString(ids)

    res, err := modules.NewUser(c.Db).RemoveUsers(uid)
    if withPost {
        res, err = modules.NewTopic(c.Db, c.Rc).RemoveUsersTopics(uid)
    }
    return c.Render(res, err)
}

func (c Admin) RemoveTopics(sid string, ids string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    tid := helpers.GetArrayFromString(ids)
    res, err    := modules.NewTopic(c.Db, c.Rc).RemoveTopics(tid)

    return c.Render(res, err)
}

func (c Admin) EditTopic(sid string, tid int64, content, keywords string) r.Result {
    session     := c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    res, err    := modules.NewTopic(c.Db, c.Rc).EditOne(session.SessionVal.UserId, tid, content, keywords, true)

    return c.Render(res, err)
}

