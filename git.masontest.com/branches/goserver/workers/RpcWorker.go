package workers

import (
    "net"
    "bufio"
    "bytes"
    "time"
    "fmt"
    "errors"
    "encoding/json"
    "git.masontest.com/branches/goserver/plugins"
    "git.masontest.com/branches/goserver/protocols"
    run "git.masontest.com/branches/goserver/runtime"
)

type RpcWorker struct {
    Id          int
    Log         *plugins.ServerLog
    Config      *map[string]string
    SecretKey   string
    TcpListener *net.TCPListener
    UdpConn     *net.UDPConn
    Proto       *protocols.TextProtocol
    Duration    time.Duration
    IsTcp       bool
    Sync        chan int
}

var (
    request *protocols.Request
    reqData *protocols.RequestData
)

func init() {
    request = &protocols.Request{}
    reqData = &protocols.RequestData{}
}

func NewRpcWorker(id int, secretKey string, config *map[string]string, listener interface{}, sync chan int) *RpcWorker {
    var worker *RpcWorker

    logger  := plugins.NewServerLog((*config)["name"])
    duration, _ := time.ParseDuration((*config)["process_timeout"] + "ms")

    if tcpListener, ok := listener.(*net.TCPListener); ok {
        worker  = &RpcWorker{
            Id:         id,
            Log:        logger,
            Config:     config,
            TcpListener:tcpListener,
            UdpConn:    nil,
            SecretKey:  secretKey,
            Proto:      protocols.NewTextProtocol(logger),
            Duration:   duration,
            IsTcp:      true,
            Sync:       sync,
        }
    } else if udpConn, ok := listener.(*net.UDPConn); ok {

        worker  = &RpcWorker{
            Id:         id,
            Log:        logger,
            Config:     config,
            TcpListener:nil,
            UdpConn:    udpConn,
            SecretKey:  secretKey,
            Proto:      protocols.NewTextProtocol(logger),
            Duration:   duration,
            IsTcp:      false,
            Sync:       sync,
        }
    }

    // currect the worker id by channel
    worker.Id, _ = <-worker.Sync
    return worker
}

func (w *RpcWorker) Run() {
    if w.IsTcp {
        w.RunTcp()
    } else {
        w.RunUdp()
    }
}

func (w *RpcWorker) RunTcp() {
    for {
        conn, err := w.TcpListener.Accept()
        if err != nil {
            continue
        }

        go w.HandleClient(conn)
    }
}

func (w *RpcWorker) RunUdp() {
    var buf []byte = make([]byte, 512)

    for {
        _, address, _ := w.UdpConn.ReadFromUDP(buf)
        w.UdpConn.SetDeadline(time.Now().Add(w.Duration))
        resp    := w.parseInput(bufio.NewReader(bytes.NewReader(buf)))

        // send response data.
        w.UdpConn.WriteToUDP(resp, address)
    }
}

func (w *RpcWorker) HandleClient(client net.Conn) {
    defer client.Close()

    // w.Log.Add("%d", (*config)["process_timeout"])
    client.SetDeadline(time.Now().Add(w.Duration))

    buf         := bufio.NewReader(client)
    resp        := w.parseInput(buf)

    client.Write(resp)
}

func (w *RpcWorker) parseInput(reader *bufio.Reader) []byte {
    resp    := &protocols.Response{
        Code: 200,
        Mesg: "",
    }

    data, err   := w.DealInput(reader)

    if err != nil {
        resp.Mesg   = fmt.Sprintf("%s", err)
        return w.Proto.Encode(resp, "RPC")
    }

    startTime   := time.Now().UnixNano()

    switch data["cmd"] {
        // at present, we just support RPC commend.
        case "RPC", "RPC:GZ":
            respData, err := w.dealRPCWork(data["data"])

            if err != nil {
                resp.Mesg   = fmt.Sprintf("%s", err)
            } else {
                resp.Data   = respData
            }

        default:
            resp.Mesg   = "RpcWorker: Oops! I am going to do nothing but RPC."
    }

    endTime     := time.Now().UnixNano()
    resp.Expend = float64(endTime - startTime) / 1.0e+09

    return w.Proto.Encode(resp, data["cmd"])
}

func (w *RpcWorker) dealRPCWork(data string) (interface{}, error) {
    // decode json data.
    err     := json.Unmarshal([]byte(data), &request)
    if err != nil { return "", err }
    err     = json.Unmarshal([]byte(request.Data), &reqData)
    if err != nil { return "", err }

    if reqData.Version != protocols.Version {
        return "", errors.New("RpcWorker: Hmm! We are now expect version " + protocols.Version)
    }

    if !w.authentication(request.Signature, request.Data) {
        return "", errors.New("RpcWorker: You want to check the RPC secret key, or the packet has broken.")
    }

    if reqData.Class=="" || reqData.Method=="" {
        return "", errors.New("RpcWorker: the class or method is empty.")
    }

    return w.DealProcess(reqData)
}

func (w *RpcWorker) authentication(signature string, data string) bool {
    return signature == w.Proto.GetActiveSignature(w.SecretKey, data)
}

func (w *RpcWorker) DealInput(reader *bufio.Reader) (map[string]string, error) {
    return w.Proto.Input(reader)
}

func (w *RpcWorker) DealProcess(req *protocols.RequestData) (interface{}, error) {
    return run.Processer((*w.Config)["name"], req)
}


