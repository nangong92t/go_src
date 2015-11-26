// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// the delete dml execute planning

package planbuilder

import (
    "sync"

    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type DeleteDML struct {
    DML
    PDML        *sqlparser.Delete
}

func NewDeleteDML(del *sqlparser.Delete, tableName string, sql []string, user *client.Client) *DeleteDML {
    d   := &DeleteDML {}
    d.SetAction("delete")

    d.PDML  = del
    d.TableName = tableName
    d.Sql   = sql
    d.User  = user
    d.Mut       = &sync.Mutex{}

    return d
}

// to come true ExecCore interface
func (d *DeleteDML) Destory() {
    d.Hosts = nil
    d.User  = nil
    d.PDML  = nil
}

func analyzeDelete(del *sqlparser.Delete, args []string, user *client.Client) (plan *ExecPlan, err error) {
	// Default plan
	plan = &ExecPlan{
		PlanId:    PLAN_PASS_DML,
		FullQuery: GenerateFullQuery(del),
	}

	tableName := sqlparser.GetTableName(del.Table)
	if tableName == "" {
		plan.Reason = REASON_TABLE
		return plan, nil
	}

    plan.Core   = NewDeleteDML(del, tableName, args, user)

	return plan, nil
}

// to come true ExecCore interface
func (d *DeleteDML) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    result, curTable, err := d.doOnAllShard()

    curTable.RowTotal   -= uint64(result.(int64))
    curTable.UpdateToRedisDB()

    return result, err
}
