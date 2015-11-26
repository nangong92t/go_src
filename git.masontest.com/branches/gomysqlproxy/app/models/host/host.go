// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package host

import (
    "fmt"
    "errors"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"

    "git.masontest.com/branches/gomysqlproxy/app/models/pooling"
)

var (
    Groups  []Group
    Hosts   []*Host
)

type Host struct {
    Id      string          // 统一生成的MD5编码id
    Ip      string
    Port    string
    User    string
    Pass    string
    Sign    string

    // ActiveMysqlConnTotal  int   // 活动mysql连接总数
    HostStatLevel   int     // 当前综合硬件评比值
    CpuFree         int     // 当前百分比
    MemFree         int     // 当前内存剩余百分比
    IoWait          int     // 当前硬盘读写等待cpu比重百分比
    HdBI            int     // 当前从硬盘每秒读入的数据量(kb)
    HdBO            int     // 当前每秒输出到硬盘的数据量(kb)
    HdFree          int     // 当前百分比
    ActiveProcess   int     // 当前活动进程数
    BlockProcess    int     // 当前挂起进程数

    IsValid         bool    // 服务器是否有效.
    Error           string

    ConnTime        int64   // all connection times

    DB              *sql.DB // 数据库连接实例
}

type Group struct {
    Id      string          // 由配置文件中给出，一旦定下后不能再更改， 否则之前数据无法再正确查询.
    Master  []Host
    Slave   []Host
}

type BestHost struct {
    Host    []*Host

    // log the current time's best host.
    LogTime int64
}

// to get the a group by group id.
func GetHostGroupById(gid string) (*Group, error) {
    if len(Groups) == 0 { return nil, errors.New("Sorry, there's no any host group") }

    for _, group := range Groups {
        if gid == group.Id { return &group, nil }
    }

    return nil, errors.New("Sorry, no found out this host group")
}

// change to get host address.
func GetHostPointer(hosts []Host) []*Host{
    res := []*Host{}

    for i:=0; i<len(hosts); i++ {
        res = append(res, &(hosts[i]))
    }

    return res
}

func (h *Host) CloseDB() (err error) {
    if h.DB != nil { 
        err = h.DB.Close()
        h.DB    = nil
    }
    return err
}

func (h *Host) ConnToDB(dbName string) (*sql.DB, error) {
    // if h.DB != nil { return h.DB, nil }

    connStrFmt  := "%s:%s@tcp([%s]:%s)/%s?charset=utf8"
    // if isWrite { connStrFmt += "&autocommit=true" }
    connStr := fmt.Sprintf(connStrFmt, h.User, h.Pass, h.Ip, h.Port, dbName)
    db, err := sql.Open("mysql", connStr)

    if err != nil {
        return nil, err
    }

    db.SetMaxIdleConns(pooling.MaxPoolingConnection)
    h.DB    = db

    return db, nil
}

func (h *Host) GetConfig(port string) map[string]string {
    if port == "" { port = h.Port }

    config  := map[string]string{
        "Uri":          fmt.Sprintf("%s:%s", h.Ip, port),
        "User":         "ignore",
        "Password":     "ignore",
        "Signature":    h.Sign,
    }

    return config
}

// To Evaluate host level, Algorithm: best = item 1 + item n:
//
//   when write:
//      item 1: conn total number * -10   -- ignore
//      item 2: conn time(int64(ns / 100000)) * -1
//      item 3: cpu free precent(000) * 5
//      item 4: bo total number(000) * -1
//      item 5: hd free precent(000) * 10
//
//   when read: item 2 .. item 3 like below, and:
//      item 6: bi total number(000) * -1
//
func (h *Host) GetHostEvaluateLevel(connTime float64, gType string) int64 {
    result  := int64(connTime / 40000) * -1
    result  += h.GetHostEvaluateLevelWithoutClient(gType)

    return result
}

// Same GetHostEvaluateLevel just without the client conn item.
func (h *Host) GetHostEvaluateLevelWithoutClient(gType string) int64 {
    result  := int64(h.CpuFree) * 5

    switch gType {
    case "write", "master":
        result  += int64(h.HdBO) * -1
        result  += int64(h.HdFree) * 10
    default:                            // default the read, slave
        result  += int64(h.HdBI) * -1
    }

    return result
}

func GetBetterHost(hosts []Host, gType string) *Host {
    return GetBetterHostByPoint(GetHostPointer(hosts), gType)
}

func GetBetterHostByPoint(hosts []*Host, gType string) *Host {
    var better *Host
    fen := int64(-999999)
    hostLen := len(hosts)

    for i:=0; i<hostLen; i++ {
        curHost := hosts[i]
        if !curHost.IsValid { continue }

        // the main hosts connetion time
        connFen := (curHost.ConnTime / 100000) * -1

        curFen  := curHost.GetHostEvaluateLevelWithoutClient(gType) + connFen
        if curFen > fen {
            fen     = curFen
            better  = curHost
        }
    }

    return better
}

// get better master group by the data size and host statu.
//
func GetBetterMasterGroup() *Group {
    maxLevel    := int64(-999999)
    var betterGroup *Group

    grpLen      := len(Groups)
    if grpLen == 1 { return &Groups[0] }

    for i:=0; i<grpLen; i++ {
        curMasterLevel  := int64(0)
        master      := Groups[i].Master
        masterLen   := len(master)

        for j:=0; j<len(master); j++ {
            curHost := &master[j]
            if !curHost.IsValid { masterLen--; continue }

            connFen := (curHost.ConnTime / 100000) * -1
            curMasterLevel  += curHost.GetHostEvaluateLevelWithoutClient("master") + connFen
        }

        curMasterLevel  = int64(curMasterLevel / int64(masterLen))

        // fmt.Printf("group %s level: %d\n", Groups[i].Id, curMasterLevel)

        if curMasterLevel > maxLevel {
            maxLevel    = curMasterLevel
            betterGroup = &Groups[i]
        }

    }

    return betterGroup
}


