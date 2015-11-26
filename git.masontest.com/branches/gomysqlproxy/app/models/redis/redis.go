// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package redis

import (
    "log"
    "time"
    "encoding/json"
    "github.com/garyburd/redigo/redis"
    "git.masontest.com/branches/gomysqlproxy/app/models/leveldb"
)

var (
    MAX_POOL_SIZE = 20

    redisPool chan redis.Conn

    DBRoot      string
    dbIp    = "123.57.52.194:6379"

    dbPool *redis.Pool

    mainKeys  = map[string]string{
        "MysqlProxy":       "___main.db___",
        "MysqlTable":       "___table.db___",
        "MysqlShardTable":  "___shardtable.db___",
        "MysqlShardDB":     "___sharddb.db___",
    }
)

func init() {
    var err error
    dbPool, err  = getRedisPool("tcp", dbIp)
    if err != nil {
        log.Fatalf("%s", err)
    }
}

// init the redis connection pool
// to use it, add a param of redis.Conn in the func declare:
func getRedisPool(proto, addr string) (*redis.Pool, error) {
    pool := &redis.Pool{
        MaxIdle:     50,
        IdleTimeout: 500,
        Dial: func() (redis.Conn, error) {
            return redis.Dial(proto, addr)
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
    _, err := pool.Get().Do("PING")
    if err != nil { return pool, err }

    return pool, nil
}

func putRedis(conn redis.Conn) {
    if redisPool == nil {
        redisPool = make(chan redis.Conn, MAX_POOL_SIZE)
    }
    if len(redisPool) >= MAX_POOL_SIZE {
        conn.Close()
        return
    }
    redisPool <- conn
}

func InitRedis(network, address string) redis.Conn {
    redisPool = make(chan redis.Conn, MAX_POOL_SIZE)
    if len(redisPool) == 0 {
        go func() {
            for i := 0; i < MAX_POOL_SIZE/2; i++ {
                c, err := redis.Dial(network, address)
                if err != nil {
                    panic(err)
                }
                putRedis(c)
            }
        }()
    }
    return <-redisPool
}

func getDbKey(as string) string {
    if as == "" { return "" }

    if name, ok := mainKeys[as]; ok {
        return name
    }

    return as
}

func DeleteDB(k, dbName string) error {
    dbKey   := getDbKey(dbName) + k

    _, err  := dbPool.Get().Do("DEL", dbKey)
    if err != nil { return err }

    return nil
}

func UpdateDB(k, v, dbName string) error {
    dbKey   := getDbKey(dbName) + k

    _, err  := dbPool.Get().Do("MSET", dbKey, v)
    if err != nil { return err }

/*
    err     =  DBConn.Flush()
    if err != nil { return err }
*/
    return nil
}

func WriteDB(k, v, dbName string) error {
    return UpdateDB(k, v, dbName)
}

func ReadDBCount(dbName, condition string) (int, error) {
    // TODO
    return 0, nil
}

func ReadDB(dbName, key string) (items []map[string]interface{}, err error) {
    dbKey   := getDbKey(dbName) + key

    reply, err := redis.Values(dbPool.Get().Do("MGET", dbKey))
    if err != nil { return nil, err }

    for _, one  := range reply {
        var data interface{}

        if one == nil { continue }

        json.Unmarshal(one.([]byte), &data)

        items   = append(items, map[string]interface{}{
            key: data.(map[string]interface{}),
        })
    }

    return items, nil
}

func BuildPrimaryKey(key string, isRand bool) string {
    return leveldb.BuildPrimaryKey(key, isRand)
}

func GetCurTime() int64 {
    return leveldb.GetCurTime()
}

func RestorePrimaryId(iids []interface{}) []string {
    return leveldb.RestorePrimaryId(iids)
}

func EncodeData(data interface{}) string {
    return leveldb.EncodeData(data)
}

func DecodeData(data string) map[string]interface{} {
    return leveldb.DecodeData(data)
}
