package modules

import (
    "std/data-service/std/app/models"
    "std/data-service/std/app/helpers"
    "github.com/coopernurse/gorp"
    "github.com/garyburd/redigo/redis"
    "html"
    "time"
    "strconv"
    "fmt"
    "errors"
    "strings"
)

type Topic struct {
    Db *gorp.DbMap
    Rc redis.Conn
}

func NewTopic(db *gorp.DbMap, rc redis.Conn) *Topic {
    return &Topic {
        Db: db,
        Rc: rc,
    }
}

func (c *Topic) AddNew(userId int64, content string, keywords string, attach *models.Attach) (*models.Topic, error) {
    // 先保存topic.
    topic   := &models.Topic {
        Id: 0,
        UserId: userId,
        Content: html.EscapeString(content),
        Created: time.Now().Unix(),
        BgId: attach.Id,
    }

    err := c.Db.Insert(topic)

    // 保存相关 keywords.
    if keywords != "" {
        c.EditKeywords(topic.Id, keywords)
    }

    if attach.Id != 0 {
        uploaderModule      := NewAttach(c.Db)
        topic.Background, _ = uploaderModule.GetAttachUrl(attach)
    }

    return topic, err
}

// save keywords.
func (m *Topic) EditKeywords(tid int64, keyword string) {
    keywordModule   := NewKeyword(m.Db)
    keywords        := keywordModule.SaveString(keyword)

    // delete old relation.
    m.Db.Exec("delete from topic_has_keyword where topic_id=?", tid)

    // 保存其关联.
    for i:=0; i<len(keywords); i++ {
        topicHasKey := &models.TopicHasKeyword{
            TopicId: tid,
            KeywordId: keywords[i].Id,
        }

        m.Db.Insert(topicHasKey)
    }
}

// Edit one topic.
func (m *Topic) EditOne(uid, tid int64, content string, keywords string, isAdmin bool) (bool, error) {
    // get old.
    var topic models.Topic
    m.Db.SelectOne(&topic, "select * from topic where topic_id=? limit 1", tid)
    if topic.Id == 0 { return false, errors.New("sorry, no found out this topic!") }
    if !isAdmin && topic.UserId != uid { return false, errors.New("sorry, you just can only edit your self's topic") }

    // update
    if content != "" { topic.Content = content }
    if keywords != "" {
        m.EditKeywords(tid, keywords)
    }

    m.Db.Update(&topic)

    return true, nil
}



// 未登录用户只能返回统一标准topic list.
func (c *Topic) GetTopicList(page int, order string, limit int, maxid int, needTotal bool) (*Pager, error) {
    var topics []models.Topic

    tableName   := "topic"

    sql     := "select * from " + tableName + " where is_del=0"
    if maxid > 0 { sql = sql + " and topic_id > " + strconv.Itoa(maxid)}

    if page == 0 { page = 1 }
    if limit == 0 { limit = 20 }

    offset  := (page - 1) * limit

    switch order {
        case "recommend":
            order = "liked_cnt desc"
        case "nearby":
            order = "topic_id desc"
        case "hottest":
            order = "ommented_cnt desc"
        default:
            order = "topic_id desc"
    }

    sql     = sql + " order by " + order
    sql     = sql + " limit " + fmt.Sprintf("%d",offset) + ", " + fmt.Sprintf("%d", limit)

    _, err  := c.Db.Select(&topics, sql)

    picIds  := []string{}
    userIds := []string{}
    topicCnt := len(topics)

    // Map 相关子id
    for i:=0; i<topicCnt; i++ {
        bgId    := fmt.Sprintf("%d", topics[i].BgId)
        userId  := fmt.Sprintf("%d", topics[i].UserId)

        if topics[i].BgId != 0 && !helpers.StringInSlice(bgId, picIds) {
            picIds = append(picIds, bgId)
        }

        if !helpers.StringInSlice(userId, userIds) {
            userIds = append(userIds, userId)
        }
    }

    total   := int64(0)
    if needTotal {
        total, _    = c.Db.SelectInt("select count(*) from " + tableName)
    }

    // 获取相关图片数据.
    var urlMap map[int64]models.Attach
    if len(picIds) > 0 {
        attachModule    := NewAttach(c.Db)
        urlMap, err     = attachModule.GetAttachMapById(picIds)
        if err != nil { return nil, err }
    }

    // 获取相关用户数据.
    var userMap map[int64]models.User
    if len(userIds) > 0 {
        userModule      := NewUser(c.Db)
        userMap, err    = userModule.GetUserMapById(userIds)
        if err != nil { return nil, err }
    }

    // Merge数据.
    for i:=0; i<topicCnt; i++ {
        topics[i].Background    = urlMap[ topics[i].BgId ].Uri
        topics[i].Author        = userMap[ topics[i].UserId ].Username
    }

    pager   := &Pager{
        Total: total,
        List: topics,
    }

    return pager, err
}

// 登陆用户可返回用户偏好性topic list.
func (c *Topic) GetTopicListWithPreferences(userId int64, page int, order string, limit int) (*Pager, error) {

    c.Rc.Do("SET", "email", "i@RincLiu.com")
    test, _ := c.Rc.Do("GET", "email")

    fmt.Printf("%s\n", test)

    topics := &Pager{
        List: 23,
    }

    return topics, nil
}

// 获取当个话题详情.
func (c *Topic) GetDetail(id int64) *models.Topic {
    var topic   models.Topic
    attachModule    := NewAttach(c.Db)
    userModule      := NewUser(c.Db)

    c.Db.SelectOne(&topic, "select * from topic where topic_id=? limit 1", id)

    // 获取背景图片.
    attach  := attachModule.GetOneAttach(topic.BgId)
    topic.Background, _ = attachModule.GetAttachUrl(attach)

    // 获取用户名
    user    := userModule.GetById(topic.UserId)
    topic.Author        = user.Username

    return &topic
}

// 获取Topic下有那些Keywords.
func (c *Topic) GetKeywordsInTopic(tid int64) []models.Keyword {
    var keywords []models.Keyword

    c.Db.Select(&keywords, "select distinct k.keyword_id, k.keyword, k.created from keyword k left join topic_has_keyword t on t.keyword_id=k.keyword_id where t.topic_id = ?", tid)

    return keywords
}

// 删除某topic.
func (c *Topic) Remove(uid int64, id int64) (bool, error) {
    topic   := c.GetDetail(id)
    if topic.UserId != uid {
        return false, errors.New("seems this topic is not belong to you.")
    }

    topic.IsDel = 1
    topic.Deleted   = time.Now().Unix()
    _, err  := c.Db.Update(topic)

    return true, err
}

// 订阅某个topic.
func (c *Topic) Subscribe(uid int64, id int64) (bool, error) {
    curSType    := models.SubscribeTypes["topic"]
    err         := NewSubscribe(c.Db).Subscribe(uid, curSType, id)

    return true, err
}

// 退订某topic.
func (c *Topic) UnSubscribe(uid int64, id int64) (bool, error) {
    curSType    := models.SubscribeTypes["topic"]
    err         := NewSubscribe(c.Db).UnSubscribe(uid, curSType, id)

    return true, err
}

// 喜欢某个topic.
func (c *Topic) Like(uid int64, id int64) (bool, error) {
    curSType    := models.LikeTypes["topic"]
    err         := NewLike(c.Db).Like(uid, curSType, id)
    if err != nil { return false, err }

    keywords    := c.GetKeywordsInTopic(id)
    NewKeyword(c.Db).ChangeUserLikeRate(uid, keywords, "+", 1)

    return true, err
}

// 不喜欢某topic.
func (c *Topic) UnLike(uid int64, id int64) (bool, error) {
    curSType    := models.LikeTypes["topic"]
    err         := NewLike(c.Db).UnLike(uid, curSType, id)

    keywords    := c.GetKeywordsInTopic(id)
    NewKeyword(c.Db).ChangeUserLikeRate(uid, keywords, "-", 5)

    return true, err
}

// 评论某Topic.
func (c *Topic) Comment(uid int64, tid int64, commentStr string) (int64, error) {
    // check the topic is exists.
    var topic models.Topic
    c.Db.SelectOne(&topic, "select * from topic where topic_id=? limit 1", tid)

    if topic.Id == 0 {
        return 0, errors.New("Sorry, the topic is not exists.")
    }

    comment     := &models.Comment{
        TopicId:    tid,
        Creator:    uid,
        ReplyTo:    0,
        Content:    html.EscapeString(commentStr),
        Created:    time.Now().Unix(),
    }
    err := c.Db.Insert(comment)
    return comment.Id, err
}

// 回复某Topic的Comment.
func (c *Topic) ReplyComment(uid int64, tid int64, commentStr string, replyToCommentId int64) (int64, error) {
    // check the topic is exists.
    var topic models.Topic
    c.Db.SelectOne(&topic, "select * from topic where topic_id=? limit 1", tid)
    if topic.Id == 0 { return 0, errors.New("Sorry, the topic is not exists.") }

    // check the comemnt is exists.
    comment     := &models.Comment{}
    c.Db.SelectOne(comment, "select * from comment where comment_id=? limit 1", replyToCommentId)
    if comment.Id == 0 { return 0, errors.New("Sorry, the reply comment is not exists.") }

    comment     = &models.Comment{
        TopicId:    tid,
        Creator:    uid,
        ReplyTo:    replyToCommentId,
        Content:    html.EscapeString(commentStr),
        Created:    time.Now().Unix(),
    }
    err := c.Db.Insert(comment)
    return comment.Id, err
}

// 获取某Topic下的comment列表.
func (c *Topic) GetCommentList(tid int64, page int, limit int) (*Pager, error) {
    comments    := []models.Comment{}
    if page == 0 { page = 1 }
    if limit == 0 { limit = 20 }
    offset := (page-1) * limit

    c.Db.Select(&comments, "select * from comment where topic_id=? order by created desc limit ?,?", tid, offset, limit)

    commentCnt  := len(comments)
    // map user ids.
    uids    := []string{}
    for i:=0; i<commentCnt; i++ {
        uid     := strconv.FormatInt(comments[i].Creator, 10)
        uids    = append(uids, uid)
    }

    // 获取相关用户数据.
    var userMap map[int64]models.User
    var err error
    if len(uids) > 0 {
        userModule      := NewUser(c.Db)
        userMap, err    = userModule.GetUserMapById(uids)
        if err != nil { return nil, err }
    }

    // Merge map.
    for i:=0; i<commentCnt; i++ {
        comments[i].Author        = userMap[ comments[i].Creator ].Username
    }

    pager   := &Pager{
        Total: 0,
        List: comments,
    }

    return pager, err
}

// Get the topic stat data in a date.
func (m *Topic) GetStat(date int64) (map[string]int64, error) {
    result  := map[string]int64 {
        "new": 0,
        "reported": 0,
        "deleted": 0,
    }

    maxDate := date + 86400
    var err error
    sqlQ    := "select count(1) from topic where is_del="
    sqlN    := " and created >= ? and created < ? limit 1"
    result["new"], err  = m.Db.SelectInt(sqlQ + "0" + sqlN, date, maxDate)
    if err != nil { return result, err }

    result["deleted"],err = m.Db.SelectInt(sqlQ + "1 and deleted >= ? and deleted < ? limit 1", date, maxDate)
    if err != nil { return result, err }

    return result, nil
}

// Get the topic comment stat data in a date.
func (m *Topic) GetCommentStat(date int64) (map[string]int64, error) {
    result  := map[string]int64 {
        "new": 0,
        "reported": 0,
        "deleted": 0,
    }

    maxDate := date + 86400
    var err error
    sqlQ    := "select count(1) from comment where is_del="
    sqlN    := " and created >= ? and created < ? limit 1"
    result["new"], err  = m.Db.SelectInt(sqlQ + "0" + sqlN, date, maxDate)
    if err != nil { return result, err }

    result["deleted"],err = m.Db.SelectInt(sqlQ + "1 and deleted >= ? and deleted < ? limit 1", date, maxDate)
    if err != nil { return result, err }

    return result, nil
}

//  get topic date news count.
func (m *Topic) GetOneMonthDateNewCnt(toDate int64) ([]interface{}, error) {
    // golang the default time is 08:00:00, so need to add to next day time
    toDate  += (24 - 8) * 3600

    days    := int64(31)
    daySeconds  := int64(86400)
    minDay  := toDate - days * daySeconds

    var res []models.DateData

    sql     := "select unix_timestamp(date_format(FROM_UNIXTIME( `created`),'%Y-%m-%d')) * 1000 as date, count(1) as val from topic where is_del=0 and created > ? and created < ? group by date"
    _, err  := m.Db.Select(&res, sql, minDay, toDate)
    if err != nil { return nil, err }

    result  := helpers.MergeDateData(res, days, toDate)

    return result, nil
}

//  get comment date news count.
func (m *Topic) GetCommOneMonthDateNewCnt(toDate int64) ([]interface{}, error) {
    // golang the default time is 08:00:00, so need to add to next day time
    toDate  += (24 - 8) * 3600

    days    := int64(31)
    daySeconds  := int64(86400)
    minDay  := toDate - days * daySeconds

    var res []models.DateData

    sql     := "select unix_timestamp(date_format(FROM_UNIXTIME(`created`),'%Y-%m-%d')) * 1000 as date, count(1) as val from comment where is_del=0 and created > ? and created < ? group by date"
    _, err  := m.Db.Select(&res, sql, minDay, toDate)
    if err != nil { return nil, err }

    result  := helpers.MergeDateData(res, days, toDate)

    return result, nil
}

// remove the user's topic.
func (m *Topic) RemoveUsersTopics(uid []interface{}) (bool, error) {
    return m.RemoveTopicsCore("user_id in (%s)", uid)
}

// remove the user's topic.
func (m *Topic) RemoveTopics(tid []interface{}) (bool, error) {
    return m.RemoveTopicsCore("topic_id in (%s)", tid)
}

// remove the topics by ids
func (m *Topic) RemoveTopicsCore(where string, val []interface{}) (bool, error) {
    valLen      := len(val)
    curTime     := time.Now().Unix()

    condition   := fmt.Sprintf(where, strings.TrimRight(strings.Repeat("?,", valLen), ","))
    sql := "update topic set is_del=1, deleted=%d where is_del=0 and %s"
    _, err  := m.Db.Exec(fmt.Sprintf(sql, curTime, condition), val...)

    return err==nil, err
}
