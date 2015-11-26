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
)

var (
    MAX_POOL_SIZE = 20

    redisPool chan redis.Conn

    DBRoot      string
    dbIp    = "123.57.52.194:6379"

    dbPool *redis.Pool
)

func init() {
    return
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

func Delete(k string) error {
    _, err  := dbPool.Get().Do("DEL", k)
    if err != nil { return err }

    return nil
}

func Update(k, v string) error {
    _, err  := dbPool.Get().Do("MSET", k, v)
    if err != nil { return err }

    return nil
}

// expried 过期时间, 单位为妙.
func Write(k, v string, expried int64) error {
    var err error
    if expried > 0 {
        _, err  = dbPool.Get().Do("SETEX", k, expried, v)
    } else {
        return Update(k, v)
    }

    if err != nil { return err }

    return nil
}

func ReadCount(dbName, condition string) (int, error) {
    // TODO
    return 0, nil
}

func Read(key string) (item interface{}, err error) {
    reply, err := redis.Values(dbPool.Get().Do("MGET", key))

    if err != nil { return nil, err }

    for _, one  := range reply {
        var data interface{}

        if one == nil { continue }

        json.Unmarshal(one.([]byte), &data)

        item    = data
        break
    }

    return item, nil
}

