// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package redis

import (
    "testing"
    "fmt"
    "code.google.com/p/go-uuid/uuid"
    "github.com/garyburd/redigo/redis"
    "time"
)

func Test_connredis(t *testing.T) {

    c := InitRedis("tcp", "123.57.52.194:6379")

    //test uuid
    startTime := time.Now()

    var Success, Failure int
    for i := 0; i < 100; i++ {
        if ok, _ := redis.Bool(c.Do("HSET", "payVerify:session", uuid.New(), "aaaa")); ok {
            Success++
            // break
        } else {
            Failure++
        }
    }
    fmt.Printf("用时：%s, 总计：100000,成功：%d, 失败：%d", time.Now().Sub(startTime), Success, Failure)
}
