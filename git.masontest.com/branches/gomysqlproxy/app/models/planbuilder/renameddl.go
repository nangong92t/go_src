// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// The Rename DDL execute plan

package planbuilder

import (
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type RenameDDL struct {
    DDL
}

func NewRenameDDL(ddl *sqlparser.DDL, sql string, user *client.Client) *RenameDDL {
    d   := &RenameDDL {}
    d.Sql   = sql
    d.PDDL  = ddl
    d.User  = user
    d.Init()

    return d
}

// To do create business
func (d *RenameDDL) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    result, curTable, err := d.alterOnAllShard("renameddl")
    if err != nil { return nil, err }

    curTable.Name   = d.NewName
    curTable.UpdateToRedisDB()

    return result, err
}
