package controllers

import (
    "std/data-service/std/app/modules"
    "std/data-service/std/app/helpers"
	r "git.masontest.com/branches/goserver/workers/revel"
    rpc "git.masontest.com/branches/goserver/client"
    "errors"
    "time"
    "fmt"
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
        case "sameip":
            list, err   = modules.Pager{}, nil
        case "sameudid":
            list, err   = modules.Pager{}, nil
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

func (c Admin) RemoveComments(sid string, ids string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    cid := helpers.GetArrayFromString(ids)
    res, err    := modules.NewTopic(c.Db, c.Rc).RemoveComment(cid)

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

func (c Admin) GetLabelList(sid string, page int, limit int) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    res, err    := modules.NewLabel(c.Db).GetList(page, limit, true)
    return c.Render(res, err)
}

func (c Admin) RemoveLabels(sid string, ids string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    lid := helpers.GetArrayFromString(ids)
    res, err    := modules.NewLabel(c.Db).RemoveLabel(lid)
    return c.Render(res, err)
}

func (c Admin) RemoveBgs(sid string, ids string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    aid := helpers.GetArrayFromString(ids)
    res, err    := modules.NewAttach(c.Db).Remove(aid)
    return c.Render(res, err)
}

func (c Admin) GetBackgroundPhotoList(sid string, page int, limit int) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    res, err    := modules.NewAttach(c.Db).GetList(page, limit, true)
    return c.Render(res, err)
}

func (c Admin) AddUnwantWord(sid string, word string) r.Result {
    session := c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    err := modules.NewUnwantWord(c.Db).Add(session.SessionVal.UserId, word)

    return c.Render(err == nil, err)
}

func (c Admin) GetUnwantWordList(sid string, page, limit int) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    res, err    := modules.NewUnwantWord(c.Db).GetList(page, limit, true)

    return c.Render(res, err)
}

func (c Admin) RemoveUnwantWord(sid string, ids string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    aid := helpers.GetArrayFromString(ids)
    res, err    := modules.NewUnwantWord(c.Db).Remove(aid)

    return c.Render(res, err)
}

func (c Admin) GetHostStatus(sid string) r.Result {
    c.Sess.GetSession(sid)
    if !c.Sess.IsAdmin() {
        return c.RenderNoAdmin()
    }

    // init the rpc client.
    rpcSecret   := r.Config.StringDefault("rpc.client.secret_key", "")
    hostStrs    := r.Config.Options("rpc.client.Service.Tracker.host")
    hosts       := []string{}
    for _, hostName := range hostStrs {
        host    := r.Config.StringDefault(hostName, "127.0.0.1:9901")
        hosts   = append(hosts, host)
    }

    conns   := map[string]*rpc.RpcClient{}
    for _, uri := range hosts {
        config := map[string]string {
            "Uri":          uri,
            "User":         "tony",
            "Password":     "123",
            "RpcSecret":    rpcSecret,
        }
        conns[uri]  = rpc.NewRpcClient("tcp", config)
    }

    result  := conns["127.0.0.1:9901"].Call("Host", "Status", nil)

    return c.Render(result, nil)

    panic(fmt.Sprintf("%#v, %#v", rpcSecret, hosts))
}
