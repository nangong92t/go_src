package models

import (
    "github.com/revel/revel"
)

type Attach struct {
    Id          int64   `db:"attach_id" json:"attach_id"`
    Size        int64   `db:"size" json:"size"`
    Extention   string  `db:"extention" json:"extention"`
    Name        string  `db:"name" json:"name"`
    Md5         string  `db:"md5" json:"md5"`
    SaveHost    string  `db:"savehost" json:"savehost"`
    SavePath    string  `db:"savepath" json:"savepath"`
    SaveName    string  `db:"savename" json:"savename"`
    Created     int64   `db:"created" json:"created"`

    Uri         string  `db:"-" json:"-"`
}

func (m *Attach) Validate(v *revel.Validation) {
    v.Check(m.Size, revel.ValidRequired())
    v.Check(m.Extention, revel.ValidRequired())
    v.Check(m.Name, revel.ValidRequired())
    // v.Check(m.Md5, revel.ValidRequired())
    v.Check(m.SavePath, revel.ValidRequired())
    v.Check(m.SaveName, revel.ValidRequired())
}

