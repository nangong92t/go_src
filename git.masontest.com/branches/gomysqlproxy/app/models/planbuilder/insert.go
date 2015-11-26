// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package planbuilder

import (
    "errors"
    "sync"

    "git.masontest.com/branches/gomysqlproxy/app/models/schema"
    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type InsertDML struct {
    DML
    PDML        *sqlparser.Insert
}

func NewInsertDML(ins *sqlparser.Insert, tableName string, sql []string, user *client.Client) *InsertDML {
    d   := &InsertDML {}
    d.SetAction("insert")

    d.PDML  = ins
    d.TableName = tableName
    d.Sql   = sql
    d.User  = user
    d.Mut       = &sync.Mutex{}

    return d
}

func analyzeInsert(ins *sqlparser.Insert, args []string, user *client.Client) (plan *ExecPlan, err error) {

	plan = &ExecPlan{
		PlanId:    PLAN_PASS_DML,
		FullQuery: GenerateFullQuery(ins),
	}
	tableName := sqlparser.GetTableName(ins.Table)

	if tableName == "" {
		plan.Reason = REASON_TABLE
		return plan, nil
	}

    plan.Core   = NewInsertDML(ins, tableName, args, user)

    return plan, nil
}

// to come true ExecCore interface
func (d *InsertDML) Destory() {
    d.Hosts = nil
    d.User  = nil
    d.PDML  = nil
}

// to come true ExecCore interface
func (d *InsertDML) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    // get the relation true table.
    var err error
    var shardTables []*schema.MysqlShardTable
    var curTable    *schema.MysqlTable
    for i:=0; i<len(schema.Tables); i++ {
        if schema.Tables[i].Name == d.TableName {
            curTable    = schema.Tables[i]
            shardTables = curTable.Shards
            break
        }
    }
    if shardTables == nil { return nil, errors.New("no found any one this table " + d.TableName) }

    // get the last true table.
    tableLen    := len(shardTables)
    var table *schema.MysqlShardTable

    for i:=0; i<tableLen; i++ {
        if shardTables[i].RowTotal < uint64(5000000) {
            table   = shardTables[i]
            break
        }
    }

    if table == nil {
        // to build new shard table
        table, err  = curTable.BuildNewShardTable()
        if err != nil { return nil, err }
    }

    // parse the sql, to find the auto increment key,
    // then change its' value to table global id.
    sql, params, err  := curTable.ParseMergeInsertGlobalId(d.Sql, table)
    if err != nil { return nil, err }

    master  := host.GetBetterHost(table.ShardDB.HostGroup.Master, "master")
    conn, err   := master.ConnToDB(table.ShardDB.Name)
    if err != nil { conn.Close(); return nil, err }

    stmt, err   := conn.Prepare(sql)
    if err != nil { conn.Close(); return nil, err }

    res, err    := stmt.Exec(params...)
    if err != nil { conn.Close(); return nil, err }

    stmt.Close()
    conn.Close()

    curTable.RowTotal++
    curTable.UpdateToRedisDB()

    return res.LastInsertId()
}


