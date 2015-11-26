// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package planbuilder

import (
    "fmt"
    "errors"
    "strings"
    "runtime"
    "database/sql"
    "sync"

    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/schema"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type DDL struct {
    Sql         string
	Action      string
	TableName   string
	NewName     string
    Hosts       []*host.Host
    PDDL        *sqlparser.DDL
    User        *client.Client

    // avoid more connection to do this op in same time.
    Mut         *sync.Mutex
}

type alterResult struct {
    Res     sql.Result
    Err     error
}

func analyzeDDL(ddl *sqlparser.DDL, sql string, user *client.Client) *ExecPlan {
	plan := &ExecPlan{PlanId: PLAN_DDL}

    switch ddl.Action {
    case "create":
        plan.Core   = NewCreateDDL(ddl, sql, user)
    case "alter":
        plan.Core   = NewAlterDDL(ddl, sql, user)
    case "rename":
        plan.Core   = NewRenameDDL(ddl, sql, user)
    case "drop":
        plan.Core   = NewDropDDL(ddl, sql, user)
    default:
        plan.Core   = NewDDL(ddl, sql, user)
    }
	// Skip TableName if table is empty (create statements) or not found in schema

	return plan
}

func NewDDL(ddl *sqlparser.DDL, sql string, user *client.Client) (plan *DDL) {
	d   := &DDL{
        Sql: sql,
        PDDL: ddl,
        User: user,
    }
    d.Init()
    return d
}

func (d *DDL) Init() {
    d.Action    = d.PDDL.Action
    d.TableName = string(d.PDDL.Table)
    d.NewName   = string(d.PDDL.NewName)
    d.Mut       = &sync.Mutex{}
}

func (d *DDL) Destory() {
    d.Hosts = nil
    d.PDDL  = nil
    d.User  = nil
    d   = nil
}

func (d *DDL) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    panic(fmt.Sprintf("%#v", d.PDDL))
    return nil, nil
}

func (d *DDL) alterOnAllShard(alterType string) (interface{}, *schema.MysqlTable, error) {
    // to found out the alter table's host.
    //var err error
    var shardTables []*schema.MysqlShardTable
    var curTable    *schema.MysqlTable
    for i:=0; i<len(schema.Tables); i++ {
        if schema.Tables[i].Name == d.TableName {
            curTable    = schema.Tables[i]
            shardTables = curTable.Shards
            break
        }
    }
    if shardTables == nil { return nil, nil, errors.New("no found any one this table " + d.TableName) }

    // when alter table name, then need change MysqlProxyTable log
    resultCh  := make(chan *alterResult, runtime.NumCPU())

    shardTBLen  := len(shardTables)

    // to do the sql at the same time.
    for i:=0; i<shardTBLen; i++ {
        go func(i int) {
            stb := shardTables[i]

            conn, err := host.GetBetterHost(stb.ShardDB.HostGroup.Master, "master").ConnToDB(stb.ShardDB.Name)
            if err != nil && i==0 {
                conn.Close()
                resultCh <- &alterResult{Res:nil, Err: err}
                return
            }

            trueSql     := ""
            if alterType == "renameddl" {
                oldSTbName  := stb.Name
                stb.Name    = strings.Replace(stb.Name, curTable.Name, d.NewName, -1)
                trueSql     = "alter table " + oldSTbName + " rename " + stb.Name
                stb.UpdateToRedisDB()
            } else if alterType == "dropddl" {
                oldSTbName  := stb.Name
                stb.Name    = strings.Replace(stb.Name, curTable.Name, d.NewName, -1)
                trueSql     = "drop table " + oldSTbName
            } else {
                trueSql     = strings.Replace(d.Sql, d.TableName, stb.Name, -1)
            }
            stmt, err   := conn.Prepare(trueSql)
            if err != nil {
                conn.Close()
                resultCh <- &alterResult{Res:nil, Err: err}
                return
            }

            res, err    := stmt.Exec()
            if err != nil {
                conn.Close()
                resultCh <- &alterResult{Res:nil, Err: err}
                return
            }

            stmt.Close()
            conn.Close()

            if i == 0 { resultCh <- &alterResult{Res:res, Err: nil}; return }
        }(i)
    }

    res := <-resultCh

    if res.Err != nil { return nil, nil, res.Err }

    result, err := res.Res.RowsAffected()

    return result, curTable, err
}
