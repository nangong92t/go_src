// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// The Go mysql proxy main model file.
package models

import (
    "fmt"
    "errors"
    "strconv"
    "sync"
    "log"

    "git.masontest.com/branches/gomysqlproxy/app/models/schema"
    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/redis"
    "git.masontest.com/branches/gomysqlproxy/app/models/planbuilder"
)

type MysqlProxy struct {
    TableTotal  uint64
    SizeTotal   uint64          // 单位：MB
    CurGId      uint64          // 当前全局表主ID，为新表主ID，不能减少
    TableIds    []string        // 表主ID
    Tables      []*schema.MysqlTable   `json:"-"`

    ShardDBIds  []string        // Shard库的主ID
    ShardDBs    []*schema.MysqlShardDB `json:"-"`  // 已生成的shard库
    ShardDBCnt  int             // shard db 计数器

    TablesMap   map[string]*schema.MysqlTable `json:"-"`
}


var (
    MyProxy  *MysqlProxy
    isUpdated  = false
)

func CheckError(err error) {
    if err != nil {
        panic(fmt.Sprintf("init the table error in mysql.go: %s", err))
    }
}

func NewMysqlProxy() *MysqlProxy {
    proxy  := &MysqlProxy {}

    proxy.Init()

    host.GetAndLogHostStatus()

    return proxy
}

// To init the necessary data.
func (m *MysqlProxy) Init() {
    m.InitMain()
    m.InitMysqlDB()
    m.InitMysqlTable()
    m.InitConnPooling()

    if isUpdated {
        // save mysql proxy.
        err     := redis.UpdateDB("main", redis.EncodeData(m), "MysqlProxy")
        CheckError(err)
    }

    // panic(fmt.Sprintf("OK: %#v", m))
}

// get the table status.
func (m *MysqlProxy) GetStatus() (map[string]interface{}, error) {
    result  := map[string]interface{}{}

    result["main"]  = redis.EncodeData(m)
    tables  := []string{}
    shardDB := []string{}

    for _, table := range m.Tables {
        tables  = append(tables, redis.EncodeData(table))
    }

    for _, db   := range m.ShardDBs {
        shardDB = append(shardDB, redis.EncodeData(db))
    }

    result["tables"]    = tables
    result["sharddbs"]  = shardDB

    return result, nil
}

// restore the main proxy data.
func (m *MysqlProxy) InitMain() {
    pr, err  := redis.ReadDB("MysqlProxy", "main")
    CheckError(err)
    if len(pr) == 0 { return }

    for _, proxy := range pr {
        proxy           = proxy["main"].(map[string]interface {})
        m.TableTotal    = uint64(proxy["TableTotal"].(float64))
        m.SizeTotal     = uint64(proxy["SizeTotal"].(float64))
        m.CurGId        = uint64(proxy["CurGId"].(float64))

        if ttableIds, isOk := proxy["TableIds"].([]interface{}); isOk && len(ttableIds) > 0 {
            m.TableIds  = redis.RestorePrimaryId(ttableIds)
        } else {
            m.TableIds  = []string{}
        }

        if dbIds, isOk  := proxy["ShardDBIds"].([]interface{}); isOk && len(dbIds) > 0 {
            m.ShardDBIds    = redis.RestorePrimaryId(dbIds)
        } else {
            m.ShardDBIds    = []string{}
        }

        m.ShardDBCnt    = int(proxy["ShardDBCnt"].(float64))
        schema.ShardDBCnt   = m.ShardDBCnt
    }

    // panic(fmt.Sprintf("%#v", m))
}

// get the current db cluster data infomations
func (m *MysqlProxy) InitMysqlDB() {
    // panic(fmt.Sprintf("%#v, %#v", m.ShardDBIds, len(m.ShardDBIds)))
    if len(m.ShardDBIds) == 0 {
        // init the shard DB
        shardDBs        := []*schema.MysqlShardDB{}
        shardDBIds      := []string{}
        m.ShardDBCnt    = 0

        for _, group    := range host.Groups {
            m.ShardDBCnt++
            shardDb, err    := m.BuildNewShardDB(&group, "shard" + strconv.Itoa(m.ShardDBCnt))

            CheckError(err)

            shardDBs    = append(shardDBs, shardDb)
            shardDBIds  = append(shardDBIds, shardDb.Id)
        }

        m.ShardDBs      = shardDBs
        m.ShardDBIds    = shardDBIds

        // to prepare save new data.
        isUpdated       = true

        // add shard dbs map.  
        schema.Sdbs     = shardDBs

    } else {
        // 分析数据，并恢复至MysqlProxy结构体中.
        shardDBs    := []*schema.MysqlShardDB{}
        for _, sid := range m.ShardDBIds {
            dbs, err := redis.ReadDB("MysqlShardDB", sid)
            CheckError(err)
            if len(dbs) != 1 { panic("no found relation shard db for id:" + sid) }

            sdb     := dbs[0][sid].(map[string]interface {})
            groupId         := sdb["HostGroupId"].(string)
            curGroup, err   := host.GetHostGroupById(groupId)
            CheckError(err)

            shardDB := &schema.MysqlShardDB{
                Id:         sdb["Id"].(string),
                Name:       sdb["Name"].(string),
                TableTotal: uint64(sdb["TableTotal"].(float64)),
                SizeTotal:  uint64(sdb["SizeTotal"].(float64)),
                HostGroupId:groupId,
                Created:    int64(sdb["Created"].(float64)),
                HostGroup:  curGroup,
            }

            shardDBs    = append(shardDBs, shardDB)
        }

        m.ShardDBs  = shardDBs

        // add shard dbs map.  
        schema.Sdbs     = shardDBs
    }

    // listen the sharddb change status.
    locker          := &sync.Mutex{}

    go func() {
        for {
            newShardDB  := <-schema.NewShardDBCh

            locker.Lock()
            defer locker.Unlock()
            m.ShardDBIds    = append(m.ShardDBIds, newShardDB.Id)
            m.ShardDBs      = append(m.ShardDBs, newShardDB)
            schema.Sdbs     = m.ShardDBs

            err     := redis.UpdateDB("main", redis.EncodeData(m), "MysqlProxy")
            if err != nil {
                log.Printf("new shard db listener error:%s", err)
            }

            m.ShardDBCnt++
            schema.ShardDBCnt   = m.ShardDBCnt
            fmt.Printf("current shard total: %d\n", schema.ShardDBCnt)

        }
    }()

    // listen the table drop action.
    go func() {
        for {
            dropedTable     := <-schema.DropedTableCh
            m.DeleteTable(dropedTable)
        }
    }()

    // panic(fmt.Sprintf("in init shard db: %#v, %#v", m))
}

func (m *MysqlProxy) DeleteTable(table *schema.MysqlTable) {
    curTables   := []*schema.MysqlTable{}
    curTableIds := []string{}

    for _, one := range m.Tables {
        if one.Name != table.Name {
            curTables   = append(curTables, one)
        }
    }

    for _, one := range m.TableIds {
        if one != table.Id {
            curTableIds   = append(curTableIds, one)
        }
    }

    // delete the relations.
    m.TableIds  = curTableIds
    m.Tables    = curTables

    err     := redis.UpdateDB("main", redis.EncodeData(m), "MysqlProxy")

    if err != nil { fmt.Printf("Delete table error when write redis: %s\n", err); return }

    schema.Tables   = curTables

    // delete selfs. 
    table.Destroy()
}

// to init or restore the table infomation.
func (m *MysqlProxy) InitMysqlTable() {
    if len(m.TableIds) == 0 { return }

    // 分析数据，并恢复至MysqlProxy结构体中.
    tables      := []*schema.MysqlTable{}

    for _, tid := range m.TableIds {
        tbs, err    := redis.ReadDB("MysqlTable", tid)
        CheckError(err)
        if len(tbs) != 1 { panic("no found relation table for id: " + tid) }
        tb          := tbs[0][tid].(map[string]interface {})
// panic(fmt.Sprintf("%#v", tbs))
        shardTbIds  := []string{}
        if std, isOk    := tb["ShardIds"].([]interface{}); isOk && len(std) > 0 {
            shardTbIds  = redis.RestorePrimaryId(std)
        }

        shardTb := []*schema.MysqlShardTable{}

        table   := &schema.MysqlTable{
            Id:         tb["Id"].(string),
            Name:       tb["Name"].(string),
            CurGId:     uint64(tb["CurGId"].(float64)),
            RowTotal:   uint64(tb["RowTotal"].(float64)),
            ShardIds:   shardTbIds,
            Created:    int64(tb["Created"].(float64)),
            Shards:     shardTb,
        }

        if len(shardTbIds) > 0 {
            // create new shard table
            shardTb, err    = m.GetShardTableByIds(shardTbIds)
            CheckError(err)
            table.Shards    = shardTb

            err     = table.RestoreColumnsByDB()
            CheckError(err)
        }

        // fmt.Printf("Init table `%s` done\n", table.Name)
        tables  = append(tables, table)
    }

    m.Tables  = tables
    schema.Tables   = m.Tables
}

// to get shard table info.
func (m *MysqlProxy) GetShardTableByIds(ids []string) ([]*schema.MysqlShardTable, error) {
    if len(ids) == 0 { return nil, nil }

    tables      := []*schema.MysqlShardTable{}
    for _, id   := range ids {
        tbs, err    := redis.ReadDB("MysqlShardTable", id)
        if err != nil { return nil, err }
        if len(tbs) != 1 { return nil, errors.New("no found the shard table for id: " + id) }

        tb      := tbs[0][id].(map[string]interface {})

        shardDbId   := tb["ShardDBId"].(string)
        shardDb,err := m.GetShardDbById(shardDbId)

        if err != nil { return nil, err }

        shardTable  := &schema.MysqlShardTable{
            Id:         tb["Id"].(string),
            Name:       tb["Name"].(string),
            RowTotal:   uint64(tb["RowTotal"].(float64)),
            ShardDBId:  shardDbId,
            Created:    int64(tb["Created"].(float64)),
            ShardDB:    shardDb,
        }

        tables      = append(tables, shardTable)
    }


    return tables, nil
}

func (m *MysqlProxy) UpdateToRedisDB() error {
    return redis.UpdateDB("main", redis.EncodeData(m), "MysqlProxy")
}

// get the shard db by ids.
func (m *MysqlProxy) GetShardDbById(sid string) (*schema.MysqlShardDB, error) {
    if sid == "" { return nil, errors.New("Sorry, the shard db id connot is empty") }

    sdb, err    := redis.ReadDB("MysqlShardDB", sid)
    if err != nil { return nil, err }

    if len(sdb) != 1 { return nil, errors.New("Load shard db wrong!") }
    tsdb        := sdb[0][sid].(map[string]interface {})
    groupId     := tsdb["HostGroupId"].(string)
    curGroup, err   := host.GetHostGroupById(groupId)
    if err != nil { return nil, err }

    shardDB     := &schema.MysqlShardDB{
        Id:         tsdb["Id"].(string),
        Name:       tsdb["Name"].(string),
        TableTotal: uint64(tsdb["TableTotal"].(float64)),
        SizeTotal:  uint64(tsdb["SizeTotal"].(float64)),
        HostGroupId:groupId,
        Created:    int64(tsdb["Created"].(float64)),
        HostGroup:  curGroup,
    }

    schema.ShardDBCnt++

    return shardDB, nil
}

// to init the connection pooling.
func (m *MysqlProxy) InitConnPooling() {
    // because the database/sql support the connection pooling
    // so just to use it.
    // 这里决定不采用预先就将所有的链接生成，还是使用到时再初始化连接.
}

func (m *MysqlProxy) BuildNewShardDB(group *host.Group, name string) (*schema.MysqlShardDB, error) {
    if name == "" { return nil, errors.New("Sorry, can not build the no name databases") }

    // init the shard db to host.
    master  := group.Master[0]
    db, err := (&master).ConnToDB("mysql")
    if err != nil { return nil, err }

    stmt, err := db.Prepare(fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci", name))
    if err != nil { return nil, err }

    _, err = stmt.Exec()
    if err != nil { return nil, err }
    stmt.Close()

    shardDbId   := redis.BuildPrimaryKey(name, true)
    shardDb := &schema.MysqlShardDB{
        Id:         shardDbId,
        Name:       name,
        TableTotal: 0,
        SizeTotal:  0,
        HostGroupId:    group.Id,
        Created:    redis.GetCurTime(),
        HostGroup:  group,
    }

    // save this new shard database to tracker.
    err     = redis.WriteDB(shardDbId, redis.EncodeData(shardDb), "MysqlShardDB")
    if err != nil { return nil, err }

    (&master).CloseDB()

    schema.ShardDBCnt++

    return shardDb, nil
}

// add a new table to mysql proxy
func (m *MysqlProxy) AddTable(tab *schema.MysqlTable) error {
    tables      := m.Tables
    tableIds    := m.TableIds

    if tables == nil {
        tables      = []*schema.MysqlTable{ tab }
        tableIds    = []string{ tab.Id }
    } else {
        tables      = append(tables, tab)
        tableIds    = append(tableIds, tab.Id)
    }
    m.Tables    = tables
    schema.Tables   = tables
    m.TableIds  = tableIds

    return m.UpdateToRedisDB()
}

// to exel the sql
func (m *MysqlProxy) Exec(args []string, user *client.Client) (interface{}, error) {
    execPlan, err   := planbuilder.GetExecPlan(args, user)
    if err != nil { return nil, err }

    result, err     := execPlan.Do()
    if err != nil { return nil, err }

    switch result.(type) {
    case *schema.MysqlTable:
        err = m.AddTable(result.(*schema.MysqlTable))
        if err == nil { result  = "create succesfully!" }
    }

    return result, err
}

// to execute exeplan
// func (m *MysqlProxy) DoExecPlan(plan *planbuilder.ExecPlan

