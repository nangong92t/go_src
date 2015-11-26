package models

import (
    "github.com/revel/revel"
)

type Like struct {
    Id          int64 `db:"like_id" json:"like_id"`
    UserId      int64 `db:"user_id" json:"user_id"`
    Type        int `db:"type" json:"type"`
    TypeId      int64 `db:"type_id" json:"type_id"`
    Created     int64 `json:"created"`
}

var LikeTypes = map[string]int{
    "topic": 1,
    "comment": 2, 
}

func (m *Like) Validate(v *revel.Validation) {
    v.Check(m.UserId, revel.ValidRequired())
    v.Check(m.Type, revel.ValidRequired())
    v.Check(m.TypeId, revel.ValidRequired())
}

