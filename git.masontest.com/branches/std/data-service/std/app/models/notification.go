package models

import (
    "git.masontest.com/branches/goserver/workers/revel"
)

type Notification struct {
    Id          int64 `db:"notification_id" json:"notification_id"`
    NoticeTo    int64 `db:"notice_to" json:"notice_to"`
    Message     string `json:"message"` 
    Created     int64 `json:"created"`
}

func (m *Notification) Validate(v *revel.Validation) {
    v.Check(m.NoticeTo, revel.ValidRequired())
    v.Check(m.Message, revel.ValidRequired())
}

