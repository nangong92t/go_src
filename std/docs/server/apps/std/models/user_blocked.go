package models

import (
    "github.com/revel/revel"
)

type UserBlocked struct {
    Id          int64 `db:"user_blocked_id" json:"-"`
    UserId      int64 `db:"user_id" json:"user_id"`
    Blocked     int64 `db:"blocked" json:"blocked"`
    Created     int64 `db:"created" json:"created"`
}

func (m *UserBlocked) Validate(v *revel.Validation) {
    v.Check(m.UserId, revel.ValidRequired())
    v.Check(m.Blocked, revel.ValidRequired())
    v.Check(m.Created, revel.ValidRequired())
}

