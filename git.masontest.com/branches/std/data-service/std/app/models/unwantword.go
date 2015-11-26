package models

import (
    "git.masontest.com/branches/goserver/workers/revel"
)

type UnwantWord struct {
    Id          int64   `db:"word_id" json:"word_id"`
    Word        string  `db:"word" json:"word"`
    Used        int     `db:"used" json:"used"`
    Creator     int64   `db:"creator" json:"creator"`
    Created     int64   `db:"created" json:"created"`

    Author      string `db:"-" json:"author"`
}

func (m *UnwantWord) Validate(v *revel.Validation) {
    v.Check(m.Word, revel.ValidRequired())
    v.Check(m.Creator, revel.ValidRequired())
    v.Check(m.Created, revel.ValidRequired())
}
