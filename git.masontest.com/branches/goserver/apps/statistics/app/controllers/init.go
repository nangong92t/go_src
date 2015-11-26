package controllers

import (
    "github.com/revel/revel"
    "strings"
    "fmt"
)

func getParamString(param string, defaultValue string) string {
    p, found := revel.Config.String(param)
    if !found {
        if defaultValue == "" {
            revel.ERROR.Fatal("Cound not find parameter: " + param)
        } else {
            return defaultValue
        }
    }

    return p
}

func getConnectionString() string {
    host := getParamString("db.host", "")
    port := getParamString("db.port", "3306")
    user := getParamString("db.user", "")
    pass := getParamString("db.password", "")
    dbname := getParamString("db.name", "auction")
    protocol := getParamString("db.protocol", "tcp")
    dbargs := getParamString("dbargs", " ")

    if strings.Trim(dbargs, " ") != "" {
        dbargs = "?" + dbargs
    } else {
        dbargs = ""
    }

    return fmt.Sprintf("%s:%s@%s([%s]:%s)/%s%s", user, pass, protocol, host, port, dbname, dbargs)
}

func getRedisConnString() (string, string) {
    host := getParamString("redis.host", "localhost")
    prot := getParamString("redis.protocol", "tcp")
    port := getParamString("redis.port", "6379")

    return prot, fmt.Sprintf("%s:%s", host, port)
}

func init() {
    // revel.OnAppStart(InitDb)
    revel.InterceptMethod((*AbstractController).Begin, revel.BEFORE)
    revel.InterceptMethod((*AbstractController).Commit, revel.AFTER)
    revel.InterceptMethod((*AbstractController).Rollback, revel.FINALLY)
}


