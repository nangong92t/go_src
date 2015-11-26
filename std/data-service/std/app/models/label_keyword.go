package models

import (
    "git.masontest.com/branches/goserver/workers/revel"
)

type LabelHasKeyword struct {
    Id          int64 `db:"label_has_keyword_id" json:"label_has_keyword_id"`
    LabelId     int64 `db:"label_id" json:"label_id"`
    KeywordId   int64 `db:"keyword_id" json:"keyword_id"`
}

func (m *LabelHasKeyword) Validate(v *revel.Validation) {
    v.Check(m.LabelId, revel.ValidRequired())
    v.Check(m.KeywordId, revel.ValidRequired())
}

