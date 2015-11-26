// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package client

import (
    "runtime"

    rpc "git.masontest.com/branches/goserver/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/host"
)

type Client struct {
    // just ip string
    IP      string

    // record the faster hosts order on here
    BestMasterHosts     *host.BestHost
    BestSlaveHosts      *host.BestHost
}

// for support to get pingto data channel 
type chanPing struct {
    Host    *host.Host
    Ping    interface{}
}

func NewClient(ip string) *Client {
    user    := &Client{
        IP: ip,
    }

    // user.OrderFasterHosts()

    return user
}

// To get best host by type
//
//  type support:
//      a: 'master' or 'write'
//      b: 'slave' or 'read'
//
func (c *Client) GetBestHost(hosts []*host.Host, hType string) *host.Host {
    hLen    := len(hosts)
    if hLen == 0 { return nil }
    if hLen == 1 { return hosts[0] }

    var bestHost *host.Host
    runChan     := make(chan *chanPing, runtime.NumCPU())

    for i:=0; i<hLen; i++ {
        go func(h *host.Host) {
            conn    := rpc.NewRpcClient("udp", h.GetConfig("9901"))
            data, err := conn.Call("host", "pingto", []interface{}{c.IP})
            if err != nil { runChan <- &chanPing{Host: h, Ping: nil}}
            runChan <- &chanPing{Host: h, Ping: data.Data}
        }(hosts[i])
    }

    bestLevel := int64(-999999)
    for i:=0; i<hLen; i++ {
        pingCh  := <-runChan
        h       := pingCh.Host
        data    := pingCh.Ping

        if res, isOk    := data.(float64); isOk {
            curLevel    := h.GetHostEvaluateLevel(res, hType)
            if curLevel > bestLevel {
                bestLevel   = curLevel
                bestHost    = h
            }
        }

    }

    close(runChan)

    return bestHost
}

// To get best host by type group
//
//  type support:
//      a: 'master'
//      b: 'slave'
//
func (c *Client) GetBestHostByGroup(group host.Group, hType string) *host.Host {
    var hosts []*host.Host
    switch hType {
    case "master":
        hosts   = host.GetHostPointer(group.Master)
    case "slave":
        hosts   = host.GetHostPointer(group.Slave)
    default:
        return nil
    }

    return c.GetBestHost(hosts, hType)
}

func (c *Client) GetCurBestMasterHost() *host.Host {

    return nil
}
