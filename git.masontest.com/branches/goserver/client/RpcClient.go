//
// To support call remote GoServer in local.
//

package client

import (
    "net"
    "time"
    "errors"
    "io/ioutil"

    "git.masontest.com/branches/goserver/plugins"
    "git.masontest.com/branches/goserver/protocols"
)

type RpcClient struct {
    // the tcp address.
    TcpAddr     *net.TCPAddr

    // the udp address.
    UdpAddr     *net.UDPAddr
    IsTcp       bool

    Log         *plugins.ServerLog
    RConfig     *protocols.RemoteConfig

    // the text protoco paraser
    TextProto   *protocols.TextProtocol

    // the default rpc commend
    Cmd         string
}

// @param connType Just support tcp or udp
// @param uri      The remote server uri
func NewRpcClient(connType string, config map[string]string) *RpcClient {
    if _, ok := config["Uri"]; !ok {
        panic("Sorry, must config the Uri for Rpc Client")
    }

    logger  :=  plugins.NewServerLog("")

    rconfig, err    := protocols.NewRemoteConfig(config)
    checkError(err, logger)

    textProto       := protocols.NewTextProtocol(logger)
    cmd             := "RPC"

    switch connType {
        case "tcp":
            tcpAddr, err := net.ResolveTCPAddr("tcp4", config["Uri"])
            checkError(err, logger)

            return &RpcClient{
                TcpAddr:    tcpAddr,
                UdpAddr:    nil,
                IsTcp:      true,
                Log:        logger,
                RConfig:    rconfig,
                TextProto:  textProto,
                Cmd:        cmd,
            }

        case "udp":
            udpAddr, err := net.ResolveUDPAddr("udp4", config["Uri"])
            checkError(err, logger)
            return &RpcClient{
                TcpAddr:    nil,
                UdpAddr:    udpAddr,
                IsTcp:      false,
                Log:        logger,
                RConfig:    rconfig,
                TextProto:  textProto,
                Cmd:        cmd,
            }
        default:
            checkError(errors.New("sorry don't support this conn type"), logger)
    }

    return nil
}

func checkError(err error, logger *plugins.ServerLog) {
    if err != nil {
        logger.Add("error in RpcClient: %s", err)
        panic(err)
    }
}

func (c *RpcClient) Call(ctrl, method string, params []interface{}) (*protocols.Response, error) {
    c.RConfig.Class     = ctrl
    c.RConfig.Method    = method
    c.RConfig.Params    = params

    packet, _   := c.TextProto.Output(c.Cmd, c.RConfig)
    result, err := c.remoteCall(packet)

    return result, err
}

// send the remote call.
func (c *RpcClient) remoteCall(packet []byte) (*protocols.Response, error) {
    if c.IsTcp {
        return c.remoteTcpCall(packet)
    } else {
        return c.remoteUdpCall(packet)
    }
}

func (c *RpcClient) remoteTcpCall(packet []byte) (*protocols.Response, error) {
    // create tcp connection
    conn, err := net.DialTCP("tcp", nil, c.TcpAddr)
    if err != nil { return nil, err }

    // send the data.
    _, err  = conn.Write(packet)
    if err != nil { return nil, err }

    resp, err   := ioutil.ReadAll(conn)
    if err != nil { return nil, err }

    if c.Cmd == "RPC:GZ" {
        resp, err = c.TextProto.Decode(resp)
        if err != nil { return nil, err }
    }

    data, err   := c.TextProto.ClientInput(resp)
    if err != nil { return nil, err }

    conn.Close()

    return data, nil
}

func (c *RpcClient) remoteUdpCall(packet []byte) (*protocols.Response, error) {
    // create udp connection
    conn, err := net.DialUDP("udp4", nil, c.UdpAddr)
    if err != nil { return nil, err }

    // send the data
    _, err = conn.Write(packet)
    if err != nil { return nil, err }

    resp := make([]byte, 4096)

    duration, _ := time.ParseDuration("1000ms")
    conn.SetDeadline(time.Now().Add(duration))
    conn.ReadFromUDP(resp)

    data, err   := c.TextProto.ClientInput(resp)
    if err != nil { return nil, err }

    conn.Close()

    return data, nil
}
