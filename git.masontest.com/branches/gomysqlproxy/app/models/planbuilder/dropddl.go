// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// The Drop DDL execute plan

package planbuilder

import (
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
    "git.masontest.com/branches/gomysqlproxy/app/models/schema"
)

type DropDDL struct {
    DDL
}

func NewDropDDL(ddl *sqlparser.DDL, sql string, user *client.Client) *DropDDL {
    d   := &DropDDL {}
    d.Sql   = sql
    d.PDDL  = ddl
    d.User  = user
    d.Init()

    return d
}

// To do create business
func (d *DropDDL) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    result, curTable, err := d.alterOnAllShard("dropddl")
    if err != nil { return nil, err }

    // notice main part to delete current table.
    schema.DropedTableCh <- curTable

    return result, nil
}
