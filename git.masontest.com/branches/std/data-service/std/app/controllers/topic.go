package controllers

import (
	r "git.masontest.com/branches/goserver/workers/revel"
    "std/data-service/std/app/modules"
    "std/data-service/std/app/models"
    "errors"
)

type Topic struct {
	AbstractController
}

// 创建话题.
func (c Topic) AddNew(sid string, content string, keywords string, bg []byte) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    attach  := &models.Attach{}
    if len(bg) > 0 {
        var err error
        // 尝试上传图片
        uploader    := modules.NewAttach(c.Db)
        attach, err = uploader.Upload(c.Request, "bg")
        if err != nil {
            return c.Render(nil, err)
        }
    }

    // 保存topic
    topicModule := modules.NewTopic(c.Db, c.Rc)
    topic, err  := topicModule.AddNew(session.SessionVal.UserId, content, keywords, attach)

	return c.Render(topic, err)
}

// add a new topic background attach.
func (c Topic) AddBackground(sid, name, extention, file string) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    result, err := modules.NewAttach(c.Db).AddNew(session.SessionVal.UserId, name, extention, file)

    return c.Render(result, err)
}

// 获取话题列表.
func (c Topic) GetTopicList(sid string, page int, order string, limit int, maxid int, needTotal bool) r.Result {
    session     := c.Sess.GetSession(sid)

    topicModule := modules.NewTopic(c.Db, c.Rc)

    var topics *modules.Pager
    var err error

    if c.Sess.IsAdmin() {
        topics, err = topicModule.GetTopicList(page, order, limit, maxid, needTotal)
    } else if session.Id == 0 {
        // 未登录用户，只能返回统一标准topic list.
        topics, err = topicModule.GetTopicList(page, order, limit, maxid, needTotal)
    } else {
        // 登陆用户可返回，用户偏好性topic list.
        topics, err = topicModule.GetTopicListWithPreferences(session.SessionVal.UserId, page, order, limit)
    }

    return c.Render(topics, err)
}

// 获取当个话题详情.
func (c Topic) GetDetail(id int64) r.Result {
    topicModule := modules.NewTopic(c.Db, c.Rc)

    topic       := topicModule.GetDetail(id)

    return c.Render(topic, nil)
}

// 订阅某个topic.
func (c Topic) Subscribe(sid string, id int64) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    if id == 0 {
        return c.Render(nil, errors.New("no topic id"))
    }

    topicModule := modules.NewTopic(c.Db, c.Rc)

    res, err    := topicModule.Subscribe(session.SessionVal.UserId, id)

    return c.Render(res, err)
}

// 退订某topic.
func (c Topic) UnSubscribe(sid string, id int64) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    if id == 0 {
        return c.Render(nil, errors.New("no topic id"))
    }

    topicModule := modules.NewTopic(c.Db, c.Rc)

    res, err    := topicModule.UnSubscribe(session.SessionVal.UserId, id)

    return c.Render(res, err)

}

// 删除某Topic.
func (c Topic) Remove(sid string, id int64) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    topicModule := modules.NewTopic(c.Db, c.Rc)
    res, err    := topicModule.Remove(session.SessionVal.UserId, id)

    return c.Render(res, err)
}


// 喜欢某Topic.
func (c Topic) Like(sid string, id int64) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    res, err    := modules.NewTopic(c.Db, c.Rc).Like(session.SessionVal.UserId, id)

    return c.Render(res, err)
}

// 不喜欢某Topic.
func (c Topic) UnLike(sid string, id int64) r.Result {
    session     := c.Sess.GetSession(sid)
    if session.Id == 0 {
        return c.RenderLogin()
    }

    res, err    := modules.NewTopic(c.Db, c.Rc).UnLike(session.SessionVal.UserId, id)

    return c.Render(res, err)
}

// 评论某Topic或回复某评论
func (c Topic) Comment(sid string, id int64, comment string, replyto int64) r.Result {
    session     := c.Sess.GetSession(sid)
    if c.Sess.IsGuest() {
        return c.RenderLogin()
    }

    if id == 0 {
        return c.Render(nil, errors.New("Sorry, no topic id."))
    }

    if comment == "" {
        return c.Render(nil, errors.New("Sorry, the comment content is empty."))
    }

    var commentId int64
    var err error

    topicModule     := modules.NewTopic(c.Db, c.Rc)
    if replyto == 0 {
        commentId, err    = topicModule.Comment(session.SessionVal.UserId, id, comment)
    } else {
        commentId, err    = topicModule.ReplyComment(session.SessionVal.UserId, id, comment, replyto)
    }

    return c.Render(map[string]int64{"comment_id":commentId}, err)
}

// 获取某Topic下的comment列表
func (c Topic) GetCommentList(id int64, page int, limit int) r.Result {
    topicModule := modules.NewTopic(c.Db, c.Rc)
    list, err   := topicModule.GetCommentList(id, page, limit, false)

    return c.Render(list, err)
}

// 获取所有comment list.
func (c Topic) GetAllCommentList(page int, limit int) r.Result {
    topicModule := modules.NewTopic(c.Db, c.Rc)
    comments, err   := topicModule.GetCommentList(0, page, limit, true)
    return c.Render(comments, err)
}


