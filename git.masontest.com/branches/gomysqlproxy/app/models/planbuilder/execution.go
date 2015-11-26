// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package planbuilder

import (
    "fmt"
    "errors"

    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

var (
    execLimit  = &sqlparser.Limit{Rowcount: sqlparser.ValArg(":_vtMaxResultSize")}
)

// ExecPlan core interface
type ExecPlanCore interface {
    Do() (interface{}, error)
    Destory()
}

// ExecPlan is to provide the right host and data mergin
//  
type ExecPlan struct {
	PlanId      PlanType
	Reason      ReasonType
	TableName   string

    // The Data put on what hosts
    Hosts       []*host.Host

	// FieldQuery is used to fetch field info
	FieldQuery  *sqlparser.ParsedQuery

	// FullQuery will be set for all plans.
	FullQuery   *sqlparser.ParsedQuery

	// For PK plans, only OuterQuery is set.
	// For SUBQUERY plans, Subquery is also set.
	// IndexUsed is set only for PLAN_SELECT_SUBQUERY
	OuterQuery  *sqlparser.ParsedQuery
	Subquery    *sqlparser.ParsedQuery
	IndexUsed   string

	// For selects, columns to be returned
	// For PLAN_INSERT_SUBQUERY, columns to be inserted
	ColumnNumbers []int

	// PLAN_PK_EQUAL, PLAN_DML_PK: where clause values
	// PLAN_PK_IN: IN clause values
	// PLAN_INSERT_PK: values clause
	PKValues    []interface{}

	// For update: set clause if pk is changing
	SecondaryPKValues []interface{}

	// For PLAN_INSERT_SUBQUERY: pk columns in the subquery result
	SubqueryPKColumns []int

	// PLAN_SET
	SetKey      string
	SetValue    interface{}

    // CORE PLAN
    Core        ExecPlanCore
}


func GetExecPlan(args []string, user *client.Client) (plan *ExecPlan, err error) {
    if len(args) == 0 { return nil, errors.New("Sorry the can not support empty sql") }
    sql         := args[0]
    if sql == "" { return nil, errors.New("Sorry the can not support empty sql") }

    // params      := args[1:]
    statement, err := sqlparser.Parse(sql)
    if err != nil { return nil, err }

	plan, err = analyzeSQL(statement, args, user)
    if err != nil { return nil, err }

	if plan.PlanId == PLAN_PASS_DML {
		fmt.Errorf("PASS_DML: %s", sql)
	}
	return plan, nil
}

func analyzeSQL(statement sqlparser.Statement, args []string, user *client.Client) (plan *ExecPlan, err error) {
	switch stmt := statement.(type) {
	case *sqlparser.DDL:
		return analyzeDDL(stmt, args[0], user), nil
    case *sqlparser.Insert:
        return analyzeInsert(stmt, args, user)
    case *sqlparser.Select:
        return analyzeSelect(stmt, args, user)
    case *sqlparser.Update:
        return analyzeUpdate(stmt, args, user)
    case *sqlparser.Delete:
        return analyzeDelete(stmt, args, user)
    case *sqlparser.Set:
        // TODO
	case *sqlparser.Other:
		return &ExecPlan{PlanId: PLAN_OTHER}, nil
	}
	return nil, errors.New("invalid SQL")
}

func (plan *ExecPlan) Do() (interface{}, error) {
    return plan.Core.Do()
}
