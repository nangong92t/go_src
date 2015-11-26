package models

import (
    "github.com/revel/revel"
)

type UserLoadedMaxTopic struct {
    UserId      int64 `db:"user_id" json:"user_id"`
    MaxTopicId  int64 `db:"max_topic_id" json:"max_topic_id"`
}

func (m *UserLoadedMaxTopic) Validate(v *revel.Validation) {
    v.Check(m.UserId, revel.ValidRequired())
    v.Check(m.MaxTopicId, revel.ValidRequired())
}

