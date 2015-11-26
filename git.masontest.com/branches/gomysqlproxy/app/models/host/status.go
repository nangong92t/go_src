// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package host

import (
    "fmt"
    "time"

    "git.masontest.com/branches/goserver/plugins"
    rpc "git.masontest.com/branches/goserver/client"
)

var (
    logger *plugins.ServerLog
)

func init() {
    logger  = plugins.NewServerLog("tony")
}

func noticeAdmin(mesg string) {
    // TODO
}

func logStatus(logType string, log interface{}) {
    // TODO: 

}

func GetAndLogHostStatus() {
    for _, host := range Hosts {
        config  := map[string]string {
            "Uri":  host.Ip+":9901",
            "User": "test",
            "Password": "123",
            "Signature": host.Sign,
        }
        conn    := rpc.NewRpcClient("udp", config)
        t1      := time.Now().UnixNano()
        data, err := conn.Call("host", "status", nil)
        t2      := time.Now().UnixNano()
        host.ConnTime   = t2 - t1
        if err != nil {
            host.Error  = fmt.Sprintf("%s", err)
            noticeAdmin(host.Error)
            host.IsValid= false
            continue
        } else if data.Code != 200 {
            host.Error = data.Mesg
            noticeAdmin(host.Error)
            host.IsValid= false
            continue
        }

        // temporarily colse
        // host.Error  = ""

        for name, one := range data.Data.(map[string]interface{}) {
            if name == "Hd" || name == "Net" {
                if oneS, isOk := one.([]interface{}); isOk {
                    switch name {
                        case "Hd":
                            parseHdState(oneS, host)
                        case "Net":
                            parseNetState(oneS, host)
                    }
                }

            } else {
                if oneS, isOk := one.(map[string]interface{}); isOk {
                    switch name {
                        case "Cpu":
                            parseCpuState(oneS, host)
                        case "Mem":
                            parseMemState(oneS, host)
                        case "Sys":
                            parseSysState(oneS, host)
                    }
                }
            }
        }

        host.IsValid= true

        logger.Add("Sys: %#v\n", host)
    }
}

func parseCpuState(cpu map[string]interface{}, host *Host) {
    // TODO: to log the cpu history to LevelDb
    // logger.Add("Cpu: %#v", cpu)
}

func parseMemState(mem map[string]interface{}, host *Host) {
    // to parse the free memery.
    total       := mem["MemTotal"].(float64)
    freeM       := mem["MemFree"].(float64)
    free        := int(freeM / total * float64(100))

    host.MemFree    = free

    // TODO: to log the mem history to LevelDb
}

func parseHdState(hds []interface{}, host *Host) {

    maxAvailable    := float64(0)
    freeP       := 0
    for _, hd   := range hds {
        tHd     := hd.(map[string]interface{})
        freeHd  := tHd["Available"].(float64)
        if freeHd > maxAvailable {
            maxAvailable    = freeHd
            totalHd := tHd["Blocks"].(float64)
            freeP   = int(freeHd / totalHd * float64(100))
        }
    }

    host.HdFree = freeP

    // TODO: to log the hd history to LevelDb
}

func parseNetState(nets []interface{}, host *Host) {
    // TODO: to log the net history to LevelDb.
    // logger.Add("Hd: %#v, %#v\n", nets, host)
}

func parseSysState(sys map[string]interface{}, host *Host) {
    // to log the free cpu percents.
    host.CpuFree        = int(sys["Id"].(float64))
    host.ActiveProcess  = int(sys["R"].(float64))
    host.BlockProcess   = int(sys["B"].(float64))
    host.IoWait         = int(sys["Wa"].(float64))
    host.HdBI           = int(sys["Bi"].(float64))
    host.HdBO           = int(sys["Bo"].(float64))

    // TODO: to log the mem history to LevelDb
}
