// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package planbuilder

import (
    "sync"
    "errors"
    "runtime"
    "strings"

    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/schema"
)

type DML struct {
    Sql         []string
	Action      string
    TableName   string
    Hosts       []*host.Host
    User        *client.Client

    // avoid more connection to do this op in same time.
    Mut         *sync.Mutex
}


func (d *DML) SetAction(action string) {
    d.Action    = action
}

// to come true ExecCore interface
func (d *DML) Destory() {
    d.Hosts = nil
    d.User  = nil
}

// to come true ExecCore interface
func (d *DML) Do() (interface{}, error) {
    return nil, nil
}

func (d *DML) doOnAllShard() (interface{}, *schema.MysqlTable, error) {
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

    // the the sql param pointer
    params      := getInterfaceParams(d.Sql[1:])

    // to do the sql at the same time.
    for i:=0; i<shardTBLen; i++ {
        go func(i int) {
            stb := shardTables[i]

            conn, err := host.GetBetterHost(stb.ShardDB.HostGroup.Master, "master").ConnToDB(stb.ShardDB.Name)
            if err != nil {
                conn.Close()
                resultCh <- &alterResult{Res:nil, Err: err}
                return
            }

            trueSql     := strings.Replace(d.Sql[0], d.TableName, stb.Name, -1)
            stmt, err   := conn.Prepare(trueSql)
            if err != nil {
                conn.Close()
                resultCh <- &alterResult{Res:nil, Err: err}
                return
            }

            res, err    := stmt.Exec(params...)
            if err != nil {
                conn.Close()
                resultCh <- &alterResult{Res:nil, Err: err}
                return
            }

            stmt.Close()
            conn.Close()

            resultCh <- &alterResult{Res:res, Err: nil}
        }(i)
    }

    result  := int64(0)
    for i:=0; i<shardTBLen; i++ {
        res := <-resultCh
        if res.Err != nil { return nil, nil, res.Err }

        rowAffected, err := res.Res.RowsAffected()
        if err != nil { return nil, nil, err }

        result  += rowAffected
    }

    return result, curTable, nil
}
