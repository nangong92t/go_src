// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package controllers

import (
    "sync"
    "github.com/tonycbcd/cron"
    "strconv"

    "git.masontest.com/branches/goserver/workers/revel"
    "git.masontest.com/branches/goserver/plugins"
    "git.masontest.com/branches/gomysqlproxy/app/models"
    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/leveldb"
    "git.masontest.com/branches/gomysqlproxy/app/models/pooling"
)

var (
    isRunning = false
)

func init() {
    revel.OnAppStart(func() {
        // 为了避免在多个app 同时运行.
        if isRunning { return }
        childProcessName    := plugins.GetArg("-server")
        if childProcessName != "GoMysqlProxy" { return }

        // loade the hosts in the config
        host.Groups   = getHostConfig()
        host.Hosts    = getAllHosts()

        // set mysql connection pooling params.
        pooling.MinPoolingConnection = revel.Config.IntDefault("mysql.ConnectionPooling.min", 10)
        pooling.MaxPoolingConnection = revel.Config.IntDefault("mysql.ConnectionPooling.max", 20)

        // init the main mysql proxy object.
        models.MyProxy = models.NewMysqlProxy()

        wg  := &sync.WaitGroup{}
        wg.Add(1)

        cron := cron.New()
        cron.AddFunc("@every 10s", func() { host.GetAndLogHostStatus() })
        cron.Start()
        isRunning   = true
    })
}

func wait(wg *sync.WaitGroup) chan bool {
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
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

func getHostConfig() []host.Group {
    groupCnt    := 0
    // output := revel.Config.StringDefault("host.group[0].master[0].ip", "fuck")
    groups      := []host.Group{}

    for {
        groupStr    := "host.group[" + strconv.Itoa(groupCnt) + "]."
        configKeys  := revel.Config.Options(groupStr)
        if len(configKeys) == 0 { break }

        masterCnt   := 0
        slaveCnt    := 0

        groupMaster := []host.Host{}
        groupSlave  := []host.Host{}

        for {
            masterStr   := groupStr + "master[" + strconv.Itoa(masterCnt) + "]."
            masterHost  := host.Host{}
            masterHost.Ip   = revel.Config.StringDefault(masterStr+"ip", "")
            if masterHost.Ip == "" { break }
            masterHost.Port = revel.Config.StringDefault(masterStr+"port", "")
            masterHost.User = revel.Config.StringDefault(masterStr+"user", "")
            masterHost.Pass = revel.Config.StringDefault(masterStr+"pass", "")
            masterHost.Sign = revel.Config.StringDefault(masterStr+"sign", "")
            masterCnt++
            masterHost.Id   = leveldb.BuildPrimaryKey(masterHost.Ip + masterHost.Port, false)
            masterHost.IsValid  = true

            groupMaster = append(groupMaster, masterHost)
        }

        for {
            slaveStr    := groupStr + "slave[" + strconv.Itoa(slaveCnt) + "]."
            slaveHost   := host.Host{}
            slaveHost.Ip    = revel.Config.StringDefault(slaveStr+"ip", "")
            if slaveHost.Ip == "" { break }

            slaveHost.Port  = revel.Config.StringDefault(slaveStr+"port", "")
            slaveHost.User  = revel.Config.StringDefault(slaveStr+"user", "")
            slaveHost.Pass  = revel.Config.StringDefault(slaveStr+"pass", "")
            slaveHost.Sign  = revel.Config.StringDefault(slaveStr+"sign", "")
            slaveCnt++
            slaveHost.Id    = leveldb.BuildPrimaryKey(slaveHost.Ip + slaveHost.Port, false)
            slaveHost.IsValid   = true

            groupSlave  = append(groupSlave, slaveHost)
        }

        if len(groupMaster) == 0 || len(groupSlave) == 0 {
            panic("the master or slave cannot empty!")
        }

        groupId := revel.Config.StringDefault(groupStr+"id", strconv.Itoa(groupCnt+1))

        groups  = append(groups, host.Group{
            Id:     groupId,
            Master: groupMaster,
            Slave:  groupSlave,
        })

        groupCnt++
    }

    return groups
}
