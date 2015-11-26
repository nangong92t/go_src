// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// The Create DDL execute plan

package planbuilder

import (
    "strings"
    "errors"

    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    //"git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/schema"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type CreateDDL struct {
    DDL
}

func NewCreateDDL(ddl *sqlparser.DDL, sql string, user *client.Client) *CreateDDL {
    d   := &CreateDDL {}
    d.Sql   = sql
    d.PDDL  = ddl
    d.User  = user
    d.Init()
    return d
}

// To do create business
func (d *CreateDDL) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    sdb     := schema.GetBetterOneShardDB("master")
    if sdb == nil {
        return nil, errors.New("Sorry, no found out any shard db in Create DDL")
    }

    if len(sdb.HostGroup.Master) == 0 {
        return nil, errors.New("Sorry, no found out any host in Create DDL")
    }

    host            := sdb.HostGroup.Master[0]

    db, err := host.ConnToDB(sdb.Name)
    if err != nil { db.Close(); return nil, err }

    oldTableName    := string(d.PDDL.NewName)
    if schema.IsExistsTable(oldTableName) {
        db.Close()
        return nil, errors.New("Sorry, the table '" + oldTableName + "' is exists!")
    }
    newTable, err   := sdb.InitTable(oldTableName, d.Sql)
    if err != nil { db.Close(); return nil, err }

    tableName   := newTable.Shards[0].Name
    d.Sql     = strings.Replace(d.Sql, oldTableName, tableName, -1)

    stmt, err   := db.Prepare(d.Sql)
    if err != nil { db.Close(); return nil, err }
    _, err = stmt.Exec()
    if err != nil { db.Close(); return nil, err }
    stmt.Close()

    db.Close()

    return newTable, nil
}
