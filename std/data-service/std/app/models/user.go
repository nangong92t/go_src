package models

import (
    "git.masontest.com/branches/goserver/workers/revel"
)

type User struct {
    Id          int64 `db:"user_id" json:"user_id"`
    Username    string `db:"username" json:"username"`
    Password    string `db:"password" json:"password"`
    Created     int64 `db:"created" json:"created"`
    UserRoleId  int64 `db:"user_role_id" json:"user_role_id"`
    IsDel       int `db:"is_del" json:"-"`
    Deleted     int64 `db:"deleted" json:"deleted"`
}

type UserProfile struct {
    Id          int64 `db:"user_id" json:"user_id"`
    Username    string `db:"username" json:"username"`
    Created     int64 `db:"created" json:"created"`

    Gender      int `db:"gender" json:"gender"`
    Age         int `db:"age" json:"age"`
    Lllness     int `db:"lllness" json:"-"`
    LllnessStr  string `db:"-" json:"lllness"`
    Avator      string `db:"avator" json:"avator"`

    Posted      int64 `db:"-" json:"posted"`
}

func (m *User) Validate(v *revel.Validation) {
    v.Check(m.Username, revel.ValidRequired(), revel.ValidMaxSize(25))
    v.Check(m.Password, revel.ValidRequired(), revel.ValidMinSize(6))
}

