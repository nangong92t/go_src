package models

import (
    "github.com/revel/revel"
)

type Label struct {
    Id          int64 `db:"label_id" json:"label_id"`
    Name        string `db:"name" json:"name"`
    Created     int64 `json:"created"`
}

func (m *Label) Validate(v *revel.Validation) {
    v.Check(m.Name, revel.ValidRequired())
}

