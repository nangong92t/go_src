package models

import (
    "github.com/revel/revel"
)

type Subscribe struct {
    Id          int64   `db:"subscribe_id" json:"subscribe_id"`
    UserId      int64   `db:"user_id" json:"user_id"`
    Type        int     `db:"type" json:"type"`
    TypeId      int64   `db:"type_id" json:"type_id"`
    Created     int64   `db:"creatd" json:"created"`
}

var SubscribeTypes  = map[string]int{
    "topic":    1,
    "user":     2,
}

func (m *Subscribe) Validate(v *revel.Validation) {
    v.Check(m.UserId, revel.ValidRequired())
    v.Check(m.Type, revel.ValidRequired())
    v.Check(m.TypeId, revel.ValidRequired())
    v.Check(m.Created, revel.ValidRequired())
}
