package workers

import (
    "net"
    "bufio"
    "time"
    "fmt"
    "errors"
    "encoding/json"
    "../plugins"
    "../protocols"
    run "../runtime"
)

type RpcWorker struct {
    Id          int
    Log         *plugins.ServerLog
    Config      *map[string]string
    Listener    *net.TCPListener
    Worker      *TextWorker
    Sync        chan int
}

func NewRpcWorker(id int, config *map[string]string, listener *net.TCPListener, sync chan int) *RpcWorker {
    logger  := plugins.NewServerLog((*config)["name"])
    worker  := &RpcWorker{
        Id:         id,
        Log:        logger,
        Config:     config,
        Listener:   listener,
        Worker:     NewTextWorker(logger),
        Sync:       sync,
    }

    // currect the worker id by channel
    worker.Id, _ = <-worker.Sync
    return worker
}

func (w *RpcWorker) Run() {
    for {
        conn, err := w.Listener.Accept()
        if err != nil {
            continue
        }

        go w.HandleClient(conn)
    }
}

func (w *RpcWorker) HandleClient(client net.Conn) {
    defer client.Close()

    resp    := &protocols.Response{
        Code: 200,
        Mesg: "",
    }

    // w.Log.Add("%d", (*config)["process_timeout"])
    duration, _ := time.ParseDuration((*w.Config)["process_timeout"] + "ms")
    client.SetDeadline(time.Now().Add(duration))

    b           := bufio.NewReader(client)
    data, err   := w.DealInput(b)

    if err != nil {
        resp.Mesg   = fmt.Sprintf("%s", err)
        client.Write(w.Worker.Encode(resp))
        return
    }

    switch data["cmd"] {
        // at present, we just support RPC commend.
        case "RPC":
            respData, err := w.dealRPCWork(data["data"])
            if err != nil {
                resp.Mesg = fmt.Sprintf("%s", err)
            } else {
                resp.Data   = respData
            }

        default:
            resp.Mesg   = "RpcWorker: Oops! I am going to do nothing but RPC."
    }

    client.Write(w.Worker.Encode(resp))
}

func (w *RpcWorker) dealRPCWork(data string) (interface{}, error) {
    // decode json data.
    request := &protocols.RequestData{}
    err     := json.Unmarshal([]byte(data), &request)
    if err != nil { return "", err }

    if request.Version != "1.0" {
        return "", errors.New("RpcWorker: Hmm! We are now expect version 1.0.")
    }

    if !w.authentication(request.Signature) {
        return "", errors.New("RpcWorker: You want to check the RPC secret key, or the packet has broken.")
    }

    if request.Class=="" || request.Method=="" {
        return "", errors.New("RpcWorker: the class or method is empty.")
    }

    return w.DealProcess(request)
}

func (w *RpcWorker) authentication(signature string) bool {
    return (*w.Config)["rpc_secret_key"] == signature
}

func (w *RpcWorker) DealInput(reader *bufio.Reader) (map[string]string, error) {
    return w.Worker.Parse(reader)
}

func (w *RpcWorker) DealProcess(req *protocols.RequestData) (interface{}, error) {
    return run.Process(w.Config["name"], req)
}
