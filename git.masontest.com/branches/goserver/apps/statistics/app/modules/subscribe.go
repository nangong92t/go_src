package modules

import (
    "github.com/coopernurse/gorp"
)

type Subscribe struct {
    Db *gorp.DbMap
}

func NewSubscribe(db *gorp.DbMap) *Subscribe {
    return &Subscribe {
        Db: db,
    }
}

func (m *Subscribe) Subscribe(uid int64, cType int, typeId int64) error {
    sql := `INSERT INTO subscribe(user_id, type, type_id, created)
        SELECT ?, ?, ?, unix_timestamp() FROM dual
        WHERE not exists (
            select * from subscribe a
            where a.user_id=? and a.type=? and a.type_id=?
        )`

    _, err      := m.Db.Exec(sql, uid, cType, typeId, uid, cType, typeId)

    return err
}

func (m *Subscribe) UnSubscribe(uid int64, cType int, typeId int64) error {
    sql := "delete from subscribe where user_id=? and type=? and type_id=?" 

    _, err      := m.Db.Exec(sql, uid, cType, typeId)

    return err
}
