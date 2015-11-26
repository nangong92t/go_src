// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// the update dml execute planning

package planbuilder

import (
    "sync"

    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type UpdateDML struct {
    DML
    PDML        *sqlparser.Update
}

func NewUpdateDML(upd *sqlparser.Update, tableName string, sql []string, user *client.Client) *UpdateDML {
    d   := &UpdateDML {}
    d.SetAction("update")

    d.PDML  = upd
    d.TableName = tableName
    d.Sql   = sql
    d.User  = user
    d.Mut       = &sync.Mutex{}

    return d
}

// to come true ExecCore interface
func (d *UpdateDML) Destory() {
    d.Hosts = nil
    d.User  = nil
    d.PDML  = nil
}

func analyzeUpdate(upd *sqlparser.Update, args []string, user *client.Client) (plan *ExecPlan, err error) {
	// Default plan
	plan = &ExecPlan{
		PlanId:    PLAN_PASS_DML,
		FullQuery: GenerateFullQuery(upd),
	}

	tableName := sqlparser.GetTableName(upd.Table)
	if tableName == "" {
		plan.Reason = REASON_TABLE
		return plan, nil
	}

    plan.Core   = NewUpdateDML(upd, tableName, args, user)

    /*
	conditions := analyzeWhere(upd.Where)
	if conditions == nil {
		plan.Reason = REASON_WHERE
		return plan, nil
	}*/

	return plan, nil
}

// to come true ExecCore interface
func (d *UpdateDML) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    result, _, err := d.doOnAllShard()

    return result, err
}
