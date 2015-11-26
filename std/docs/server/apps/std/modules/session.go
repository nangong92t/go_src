package modules

import (
    "crypto/md5"
    "github.com/coopernurse/gorp"
    "std/data-service/std/app/models"
    "encoding/json"
    "time"
    "fmt"
)

type Session struct {
    Db *gorp.DbMap
    Session *models.Session
}

func NewSession(db *gorp.DbMap) *Session {
    return &Session {
        Db: db,
        Session: nil,
    }
}

// 创建Session
func (m *Session) Build(user *models.User) string {
    // 清除所有已过期Session.
    m.cleanExpried()

    // 如果存在当前用户未过期session, 则直接返回
    session := m.GetSessionByUserId(user.Id)
    if session.Id != 0 {
        return session.Session
    }

    // 生成新的Session.
    sessionByte := []byte(user.Username + fmt.Sprintf("%x", time.Now().UnixNano() / int64(time.Millisecond)) + user.Password)
    sessionStr := fmt.Sprintf("%x", md5.Sum(sessionByte))
    sessionVal := models.SessionVal {
        UserId: user.Id,
        Name: user.Username,
        RoleId: user.UserRoleId,
    }

    sessionValStr, _ := json.Marshal(sessionVal)
    nowTime := time.Now().Unix()

    session = &models.Session {
        Id: 0,
        Session: sessionStr,
        UserId: user.Id,
        Val: fmt.Sprintf("%s", sessionValStr),
        Created: nowTime,
        Expried: nowTime + 86400,               // 默认为1天
    }

    m.Db.Insert(session)

    return sessionStr
}

// 删除所有过期Session.
func (m *Session) cleanExpried() {
    m.Db.Exec("delete from user_session where expried < unix_timestamp()")
}

// 获取已当前用户已存在的Session.
func (m *Session) GetSessionByUserId(userId int64) *models.Session {
    var session models.Session

    m.Db.SelectOne(&session, "select session_id, session_key from user_session where user_id=?", userId)
    return &session
}

// get user infomation by session id.
func (m *Session) GetSession(sid string) *models.Session {

    var sessionVal models.SessionVal
    var session models.Session

    if sid == "" {
        return &session
    }

    m.Db.SelectOne(&session, "select * from user_session where session_key=?", sid)

    if session.Id == 0 {
        return &session
    }

    if session.Val != "" {
        json.Unmarshal([]byte(session.Val), &sessionVal)
        session.SessionVal  = &sessionVal
    }

    m.Session   = &session

    return &session
}

// 删除某Session
func (m *Session) Remove(sessionKey string) bool {
    _, err := m.Db.Exec("delete from user_session where session_key = ?", sessionKey)

    if err != nil {
        return false
    } else {
        return true
    }
}

// 检测当前用户是否已登陆.
func (m *Session) IsGuest() bool {
    return m.Session == nil || m.Session.Id == 0
}

// 检测当前用户是否有管理权限.
func (m *Session) IsAdmin() bool {
    return m.Session != nil && m.Session.SessionVal.RoleId == 1
}
