package modules

import (
    "github.com/coopernurse/gorp"
)

type Like struct {
    Db *gorp.DbMap
}

func NewLike(db *gorp.DbMap) *Like {
    return &Like{
        Db: db,
    }
}

func (m *Like) Like(uid int64, cType int, typeId int64) error {
    sql := `INSERT INTO` + "`like`" + `(user_id, type, type_id, created)
        SELECT ?, ?, ?, unix_timestamp() FROM dual
        WHERE not exists (
            select * from `+"`like`" +` a
            where a.user_id=? and a.type=? and a.type_id=?
        )`

    _, err      := m.Db.Exec(sql, uid, cType, typeId, uid, cType, typeId)

    return err
}

func (m *Like) UnLike(uid int64, cType int, typeId int64) error {
    sql := "delete from `like` where user_id=? and type=? and type_id=?" 

    _, err      := m.Db.Exec(sql, uid, cType, typeId)

    return err
}
