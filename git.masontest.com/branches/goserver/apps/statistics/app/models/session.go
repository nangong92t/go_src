package models

import (
    "github.com/revel/revel"
)

type SessionVal struct {
    UserId int64
    Name string
    RoleId int64
}

type Session struct {
    Id          int64   `db:"session_id" json:"session_id"`
    Session     string  `db:"session_key" json:"session"`
    Val         string  `db:"session_val" json:"val"`
    UserId      int64   `db:"user_id"`
    Created     int64   `db:"created" json:"created"`
    Expried     int64   `db:"expried" json:"expried"`
    SessionVal  *SessionVal `db:"-" json:"-"`
}

func (m *Session) Validate(v *revel.Validation) {
    v.Check(m.Session, revel.ValidRequired())
}

