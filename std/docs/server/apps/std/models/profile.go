package models

import (
    "github.com/revel/revel"
)

type Profile struct {
    Id          int64 `db:"user_id" json:"user_id"`
    Gender      int `db:"gender" json:"gender"`
    Age         int `db:"age" json:"age"`
    Lllness     int `db:"lllness" json:"lllness"`
    Avator      string `db:"avator" json:"avator"`
}

var (
    LllnessTypes    = map[int]string{
        1:      "HSV-1 (herpes type 1, usually cold sore)",
        2:      "HSV-1 (herpes type 1, usually genital)",
        4:      "HSV-2 (herpes type 2, usually genital)",
        8:      "Herpes (not sure which type)",
        16:     "HPV (human papillomavirus)",
        32:     "HIV (human immunodeficiency virus)",
        64:     "Hepatitis B",
        128:    "Hepatitis C",
        256:    "Chlamydia",
        512:    "Thrush",
        1024:   "Syphilis",
        2048:   "Gonorrhea",
        4096:   "Other (not on list)",
    }

    GenderTypes     = map[int]string {
        1:      "man",
        2:      "woman",
    }

)

func (m *Profile) Validate(v *revel.Validation) {
    v.Check(m.Id, revel.ValidRequired())
    v.Check(m.Gender, revel.ValidRequired())
    v.Check(m.Age, revel.ValidRequired())
    v.Check(m.Lllness, revel.ValidRequired())
}

