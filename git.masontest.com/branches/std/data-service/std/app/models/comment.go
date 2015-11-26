package models

import (
    "git.masontest.com/branches/goserver/workers/revel"
)

type Comment struct {
    Id          int64 `db:"comment_id" json:"comment_id"`
    TopicId     int64 `db:"topic_id" json:"topic_id"`
    Creator     int64 `db:"creator" json:"creator"`
    ReplyTo     int64 `db:"reply_to" json:"reply_to"`
    Content     string `db:"content" json:"content"`
    Created     int64 `json:"created"`

    Author      string `db:"-" json:"author"`
}

func (m *Comment) Validate(v *revel.Validation) {
    v.Check(m.TopicId, revel.ValidRequired())
    v.Check(m.Creator, revel.ValidRequired())
    v.Check(m.Content, revel.ValidRequired())
}

