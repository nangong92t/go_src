// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package schema

import (
    "time"
    "fmt"
    "git.masontest.com/branches/gomysqlproxy/app/models/redis"
    "git.masontest.com/branches/gomysqlproxy/app/models/host"
)

const MaxTableTotalInOneDB  = uint64(512)

type MysqlShardDB struct {
    Id          string
    Name        string                  // shard库名
    TableTotal  uint64                  // 当前库中的总表数
    SizeTotal   uint64                  // 当前库数据总量， 单位：MB
    HostGroupId string                  // 方便服务启动后能正确从leveldb中恢复相关数据到内存.
    Created     int64

    HostGroup   *host.Group      `json:"-"`  // 当前Shard库所存放在那个集群group中, 不会存储到数据库中.
}

type MysqlShardTable struct {
    Id          string
    Name        string                      // shard表名
    RowTotal    uint64                      // 当前shard表中数据总条数.
    ShardDBId   string
    Created     int64

    ShardDB     *MysqlShardDB   `json:"-"`  // 存放在那个shard DB中.
}

type MysqlTable struct {
    Id          string
    Name        string
    CurGId      uint64          // 当前表全局ID, 默认为1.
    RowTotal    uint64          // 表中数据总条数
    Created     int64           // 创建时间
    ShardIds    []string        // Shard表的主ID
    Shards      []*MysqlShardTable  `json:"-"`

    Columns     []TableColumn   `json:"-"`
    Indexes     []*Index         `json:"-"`
    PKColumns   []int
	CacheType   int
}

func (sdb *MysqlShardDB) InitTable(tName string, sql string) (*MysqlTable, error) {
    // 为了避免高并发时，整体tabletotal延迟执行导致的total数据不准确，
    // 所有在后面如果出现错误时， 需要将此添加回退.
    sdb.TableTotal  += 1

    tableId     := redis.BuildPrimaryKey(tName, true)
    newTable    := &MysqlTable {
        Id:         tableId,
        Name:       tName,
        CurGId:     1,
        RowTotal:   0,
        Created:    time.Now().Unix(),
		// Columns: make([]TableColumn, 0, 16),
		Indexes: make([]*Index, 0, 8),
    }

    newShardTable, err   := sdb.InitShardTable(tName)
    if err != nil { sdb.TableTotal -= 1; return nil, err }

    newTable.ShardIds   = []string{newShardTable.Id}
    newTable.Shards     = []*MysqlShardTable{ newShardTable }

    // add column data to memery.
    newTable.RestoreColumnsByDB()

    err     = redis.WriteDB(tableId, redis.EncodeData(newTable), "MysqlTable")
    if err != nil { sdb.TableTotal -= 1; return nil, err }

    sdb.UpdateToRedisDB()

    return newTable, nil
}

func (sdb *MysqlShardDB) InitShardTable(tName string) (*MysqlShardTable, error) {
    shardTabId  := redis.BuildPrimaryKey(tName, true)
    newSTable   := &MysqlShardTable {
        Id:         shardTabId,
        Name:       tName + "_shard1",
        RowTotal:   0,
        ShardDBId:  sdb.Id,
        Created:    time.Now().Unix(),
        ShardDB:    sdb,
    }

    err     := redis.WriteDB(shardTabId, redis.EncodeData(newSTable), "MysqlShardTable")
    if err != nil { return nil, err }

    return newSTable, nil
}

func (sdb *MysqlShardDB) UpdateToRedisDB() error {
    return redis.UpdateDB(sdb.Id, redis.EncodeData(sdb), "MysqlShardDB")
}

func (stb *MysqlShardTable) UpdateToRedisDB() error {
    return redis.UpdateDB(stb.Id, redis.EncodeData(stb), "MysqlShardTable")
}

// get one better shard db in more shard dbs
//
// @hType can support the master ans slave
//
func GetBetterOneShardDB(hType string) *MysqlShardDB {
    if len(Sdbs) == 0 { return nil }
    maxLevel    := int64(-999999)
    var betterShardDb *MysqlShardDB

    for i:=0; i<len(Sdbs); i++ {
        curGroupLevel   := int64(0)
        hosts   := Sdbs[i].HostGroup.Master
        if hType == "slave" {
            hosts =  Sdbs[i].HostGroup.Slave
        }

        for j:=0; j<len(hosts); j++ {
            curHost         := hosts[j]
            curGroupLevel   += curHost.GetHostEvaluateLevel(float64(curHost.ConnTime), hType)
        }

        if curGroupLevel >= maxLevel {
            maxLevel    = curGroupLevel
            betterShardDb   = Sdbs[i]
        }
    }

    if hType == "master" && (betterShardDb == nil || betterShardDb.TableTotal > MaxTableTotalInOneDB) {
        var err error
        betterShardDb, err   = buildNewShardDB(betterShardDb.HostGroup)
        if err != nil {
            fmt.Printf("build new sharddb error: %s\n", err)
        }
    }

    return betterShardDb
}

