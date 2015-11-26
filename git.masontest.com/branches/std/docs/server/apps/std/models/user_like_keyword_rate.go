package models

import (
    "github.com/revel/revel"
)

type UserLikeKeywordRate struct {
    UserId      int64 `db:"user_id" json:"user_id"`
    KeywordId   int64 `db:"keyword_id" json:"keyword_id"`
    Rate        int `db:"rate" json:"rate"`
}

func (m *UserLikeKeywordRate) Validate(v *revel.Validation) {
    v.Check(m.UserId, revel.ValidRequired())
    v.Check(m.KeywordId, revel.ValidRequired())
    v.Check(m.Rate, revel.ValidRequired())
}

