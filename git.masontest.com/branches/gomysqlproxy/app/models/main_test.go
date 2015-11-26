// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package models

import (
    "testing"
    "fmt"

    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/pooling"
)

var M *MysqlProxy

func init() {
    host.Groups   = []host.Group {
        {
            Id:     "1",
            Master: []host.Host {
                {
                    Ip:     "123.57.37.140",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
                    CpuFree: 100,
                    HdBI:   1,
                    HdBO:   4,
                    HdFree: 83,
                    IsValid: true,
                    ConnTime:28023751,
                },
            },
            Slave:  []host.Host {
                {
                    Ip:     "112.124.114.237",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
                    CpuFree: 100,
                    HdBI:   1,
                    HdBO:   3,
                    HdFree: 84,
                    IsValid: true,
                    ConnTime:6413248,
                },
                {
                    Ip:     "182.92.150.51",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
                    CpuFree: 97,
                    HdBI:   58,
                    HdBO:   382,
                    HdFree: 68,
                    IsValid: true,
                    ConnTime:133486911,
                },
            },
        },
        {
            Id:     "2",
            Master: []host.Host {
                {
                    Ip:     "121.199.44.61",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
                    CpuFree: 92,
                    HdBI:   0,
                    HdBO:   17,
                    HdFree: 37,
                    IsValid: true,
                    ConnTime:5206242,
                },
            },
            Slave:  []host.Host {
                {
                    Ip:     "114.215.184.73",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
                    CpuFree: 100,
                    HdBI:   2,
                    HdBO:   4,
                    HdFree: 57,
                    IsValid: true,
                    ConnTime:6303456,
                },
                {
                    Ip:     "182.92.130.30",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
                    CpuFree: 100,
                    HdBI:   2,
                    HdBO:   4,
                    HdFree: 57,
                    IsValid: true,
                    ConnTime:29930916,
                },
            },
        },

    }

    host.Hosts    = getAllHosts()

    pooling.MinPoolingConnection = 10
    pooling.MaxPoolingConnection = 20
}

func getAllHosts() []*host.Host {
    hosts := []*host.Host{}

    // get all hosts.
    for _, config := range host.Groups {
        for i:=0; i<len(config.Master); i++ {
            hosts   = append(hosts, &config.Master[i])
        }

        for i:=0; i<len(config.Slave); i++ {
            hosts   = append(hosts, &config.Slave[i])
        }
    }

    return hosts
}

func Test_InitMysqlTable(t *testing.T) {
    M = NewMysqlProxy()
    if len(M.ShardDBs) != 2 {
        // t.Errorf("sorry , the proxy init error")
    }
}

func Test_ExecSql(t *testing.T) {
    user    := client.NewClient("58.96.178.236")
    // res, err   := M.Exec("select * from test where id=? and created>?", "28", "45")
    // res, err   := M.Exec([]string{"create table maybe(id int(11) not null primary key auto_increment, name varchar(40))"}, user)
    // res, err   := M.Exec([]string{"insert into maybe values('', ?)", "TonyXu"}, user)
    // res, err   := M.Exec([]string{"alter table maybe1 change name name1 varchar(20) not null"}, user)
    // res, err   := M.Exec([]string{"alter table maybe2 rename maybe"}, user)
    // res, err   := M.Exec([]string{"create index test_name on maybe(name1)"}, user)
    res, err   := M.Exec([]string{"drop table maybe"}, user)
    // res, err   := M.Exec([]string{"select * from maybe"}, user)
    //res, err   := M.Exec([]string{"select * from maybe limit 1,2"}, user)
    //res, err   := M.Exec([]string{"select * from maybe limit 1, 2"}, user)
    // res, err   := M.Exec([]string{"select * from maybe order by id desc limit 1, 2"}, user)
    //res, err   := M.Exec([]string{"select * from maybe where name1=?", "TonyXu"}, user)
    // res, err   := M.Exec([]string{"update maybe set name1=? where id=?", "TonyHahaha", "1"}, user)
    // res, err   := M.Exec([]string{"delete from maybe where id=?", "1"}, user)
    fmt.Printf("parser result:%#v, %#v\n", res, err)
}
