// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package client

import (
    "fmt"
    "testing"

    "git.masontest.com/branches/gomysqlproxy/app/models/host"
)

var C *Client

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
                },
                {
                    Ip:     "121.199.44.61",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
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
                },
                {
                    Ip:     "114.215.184.73",
                    Port:   "3306",
                    User:   "root",
                    Pass:   "shaluo",
                    Sign:   "ab1f8e61026a7456289c550cb0cf77cda44302b4",
                    CpuFree: 98,
                    HdBI:   2,
                    HdBO:   5,
                    HdFree: 94,
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
                },
            },
        },
    }

    host.Hosts    = getAllHosts()
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

func Test_Init(t *testing.T) {
    C   = NewClient("58.96.178.236")
}

func Test_GetFastMasterHost(t *testing.T) {
    fasterOne   := C.GetBestHostByGroup(host.Groups[0], "master")
    if fasterOne.Ip != "123.57.37.140" {
        t.Errorf("check Error!")
    }
}

func Test_GetFastSlaveHost(t *testing.T) {
    fasterOne   := C.GetBestHostByGroup(host.Groups[0], "slave")
    fmt.Printf("%#v\n", fasterOne)
}
