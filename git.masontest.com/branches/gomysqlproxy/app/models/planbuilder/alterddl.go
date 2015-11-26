// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// The Create DDL execute plan

package planbuilder

import (
    "strconv"

    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqltypes"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type AlterDDL struct {
    DDL
}

func NewAlterDDL(ddl *sqlparser.DDL, sql string, user *client.Client) *AlterDDL {
    d   := &AlterDDL {}
    d.Sql   = sql
    d.PDDL  = ddl
    d.User  = user
    d.Init()

    return d
}

// To do create business
func (d *AlterDDL) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    result, curTable, err := d.alterOnAllShard("alterddl")
    if err != nil { return nil, err }

    // parse the sql.
    cols    := sqlparser.ParseAlterPart(d.Sql, d.TableName)
    for colName, options  := range cols {
        alterAction,_ := strconv.Atoi(options["action_id"])
        colExtra    := ""
        if _, isOk := options["auto_increment"]; isOk {
            colExtra    = "auto_increment"
        }

        switch alterAction {
        case sqlparser.ADD_ALTER:
            curTable.AddColumn(colName, options["type"], sqltypes.MakeString([]byte("")), colExtra)
        case sqlparser.DROP_ALTER:
            curTable.RemoveColumn(colName)
        case sqlparser.CHANGE_ALTER:
            curTable.ChangeColumn(colName, options)
        }
    }

    return result, nil
}
