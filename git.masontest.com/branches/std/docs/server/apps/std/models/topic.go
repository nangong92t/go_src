package models

import (
    "github.com/revel/revel"
)

type Topic struct {
    Id          int64 `db:"topic_id" json:"topic_id"`
    UserId      int64 `db:"user_id" json:"user_id"`
    Content     string `db:"content" json:"content"`
    Created     int64 `json:"created"`
    BgId        int64 `db:"bgid" json:"-"`
    Commented   int64 `db:"commented_cnt" json:"commented_cnt"`
    Liked       int64 `db:"liked_cnt" json:"liked_cnt"`
    IsDel       int   `db:"is_del" json:"-"`
    Deleted     int64 `db:"deleted" json:"-"`

    Background  string  `db:"-" json:"background"`
    Author      string `db:"-" json:"author"`
}

func (m *Topic) Validate(v *revel.Validation) {
    v.Check(m.UserId, revel.ValidRequired())
    v.Check(m.Content, revel.ValidRequired())
}

