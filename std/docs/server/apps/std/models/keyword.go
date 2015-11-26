package models

import (
    "github.com/revel/revel"
)

type Keyword struct {
    Id          int64 `db:"keyword_id" json:"keyword_id"`
    Keyword     string `db:"keyword" json:"keyword"`
    Created     int64 `json:"created"`
}

func (m *Keyword) Validate(v *revel.Validation) {
    v.Check(m.Keyword, revel.ValidRequired())
}

