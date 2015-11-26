package models

import (
    "github.com/revel/revel"
)

type TopicHasKeyword struct {
    TopicId     int64 `db:"topic_id" json:"topic_id"`
    KeywordId   int64 `db:"keyword_id" json:"keyword_id"`
}

func (m *TopicHasKeyword) Validate(v *revel.Validation) {
    v.Check(m.TopicId, revel.ValidRequired())
    v.Check(m.KeywordId, revel.ValidRequired())
}

