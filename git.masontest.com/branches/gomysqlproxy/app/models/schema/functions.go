// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// include more core features.

package schema

import (
    "fmt"
    "log"
    "errors"
    "time"
    "strconv"
    "strings"
    "database/sql"
    "sync"

    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/redis"
)

var (
    NewShardDBCh    = make(chan *MysqlShardDB)
    DropedTableCh   = make(chan *MysqlTable)
    ShardDBCnt      = 0
    Tables          []*MysqlTable
    Sdbs            []*MysqlShardDB
    mut             = &sync.Mutex{}
)

// get the better slave host db connection by shard tables. 
//
func (tbl *MysqlTable) GetSlaveShardDBConn(shardTabOrderId int) (*sql.DB, error){

    if len(tbl.Shards) == 0 { return nil, errors.New("No any one shard table") }

    shardDB := tbl.Shards[shardTabOrderId].ShardDB

    h       := host.GetBetterHost(shardDB.HostGroup.Slave, "slave")

    if h == nil { return nil, errors.New("no found out any one valid host") }

    db, err := h.ConnToDB(shardDB.Name)
    if err != nil { return nil, err }

    return db, err
}

// destroy self.
func (tbl *MysqlTable) Destroy() {
    redis.DeleteDB(tbl.Id, "MysqlTable")
    tbl = nil
}

// get the shard db. 
//
// @param shardDBCnt is max shard database total number.
//
func (tbl *MysqlTable) GetMasterShardDBByGroup(grp *host.Group) (*MysqlShardDB, error) {
    shardsLen   := len(tbl.Shards)

    if shardsLen > 0 && tbl.Shards[shardsLen-1].ShardDB.TableTotal < MaxTableTotalInOneDB {
        return tbl.Shards[shardsLen-1].ShardDB, nil
    }

    return buildNewShardDB(grp)
}

// build a new shard DB after more the limit
// or no any one shard DB.
func buildNewShardDB(grp *host.Group) (*MysqlShardDB, error) {
    mut.Lock()
    defer mut.Unlock()

    shardDBName := "shard" + strconv.Itoa(ShardDBCnt + 1)

    // to check this new shard db has been exists.
    isExists    := IsExistsShardDB(shardDBName)
    if isExists != nil { return isExists, nil }

    newShardDBId:= redis.BuildPrimaryKey(shardDBName, true)

    shardDB     := &MysqlShardDB{
        Id:         newShardDBId,
        Name:       shardDBName,
        TableTotal: uint64(0),
        SizeTotal:  uint64(0),
        HostGroupId:grp.Id,
        Created:    time.Now().Unix(),
        HostGroup:  grp,
    }

    // create the database to host.
    //
    master  := host.GetBetterHost(grp.Master, "master")
    db, err := (&master).ConnToDB("mysql")
    if err != nil { return nil, err }

    stmt, err := db.Prepare(fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci", shardDBName))
    if err != nil { return nil, err }

    _, err = stmt.Exec()
    if err != nil { return nil, err }
    stmt.Close()
    db.Close()

    // to write the new shard db to redis 
    //
    err     = redis.WriteDB(newShardDBId, redis.EncodeData(shardDB), "MysqlShardDB")
    if err != nil { return nil, err }

    // notice Mysql Project object to add a new shard db memery.
    NewShardDBCh <- shardDB

    return shardDB, nil
}

// restore the schema ddl.
func (tbl *MysqlTable) BuildNewShardTable() (*MysqlShardTable, error) {
    curSharded  := len(tbl.Shards)
    tTableName  := tbl.Name + "_shard" + strconv.Itoa(curSharded + 1)

    newTabId  := redis.BuildPrimaryKey(tTableName, true)

    shardTable  := &MysqlShardTable{
        Id:         newTabId,
        Name:       tTableName,
        RowTotal:   uint64(0),
        ShardDBId:  "",
        Created:    time.Now().Unix(),
        ShardDB:    nil,
    }

    // get bettet group.
    group       := host.GetBetterMasterGroup()

    shardDB, err    := tbl.GetMasterShardDBByGroup(group)
    if err != nil { return nil, err }

    shardTable.ShardDBId    = shardDB.Id
    shardTable.ShardDB      = shardDB

    betterHost  := host.GetBetterHost(group.Master, "master")

    ddlSql, err := tbl.GetSchemaDDLByDb()
    if err != nil { return nil, err }

    ddlSql  = strings.Replace(ddlSql, tbl.Shards[0].Name, tTableName, -1)
    db, err := betterHost.ConnToDB(shardDB.Name)
    if err != nil { return nil, err }

    stmt, err   := db.Prepare(ddlSql)
    if err != nil { return nil, err }

    _, err  = stmt.Exec()
    if err != nil { return nil, err }

    stmt.Close()

    // write the new shard table to redis 
    err     = redis.WriteDB(newTabId, redis.EncodeData(shardTable), "MysqlShardTable")
    if err != nil { return nil, err }

    tbl.Shards      = append(tbl.Shards, shardTable)
    tbl.ShardIds    = append(tbl.ShardIds, newTabId)

    err     = redis.UpdateDB(tbl.Id, redis.EncodeData(tbl), "MysqlTable")
    if err != nil { return nil, err }

    return shardTable, nil
}

func (tbl *MysqlTable) UpdateToRedisDB() error {
    return redis.UpdateDB(tbl.Id, redis.EncodeData(tbl), "MysqlTable")
}

// to get the ddl sql
//
func (tbl *MysqlTable) GetSchemaDDLByDb() (string, error) {
    db, err := tbl.GetSlaveShardDBConn(0)
    if err != nil { return "", err }

    // var rows *sql.Rows
    rows, err := db.Query(fmt.Sprintf("show create table `%s`", tbl.Shards[0].Name))
    if err != nil { return "", err }

    ddl     := ""
    if rows != nil {
        for rows.Next() {
            var tableName string

            err     = rows.Scan(&tableName, &ddl)
            if err != nil { return "", err }
        }

        rows.Close()
    }

    return ddl, nil
}

func (tbl *MysqlTable) GetGId() string {
    curId   := tbl.CurGId
    tbl.CurGId++

    err     := redis.UpdateDB(tbl.Id, redis.EncodeData(tbl), "MysqlTable")
    if err != nil {
        log.Printf("get gid error: %s", err)
    }

    return strconv.FormatUint(curId, 10)
}

func IsExistsTable(tableName string) bool {
    isFound := false
    if len(Tables) == 0 { return isFound }

    for _, one := range Tables {
        if one.Name == tableName {
            isFound = true
            break
        }
    }

    return isFound
}

func IsExistsShardDB(shardDbName string) *MysqlShardDB {
    var isExists *MysqlShardDB

    for _, one  := range Sdbs {
        if one.Name == shardDbName {
            isExists    = one
            break
        }
    }

    return isExists
}
