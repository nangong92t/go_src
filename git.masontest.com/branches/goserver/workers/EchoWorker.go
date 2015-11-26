package workers

import (
    "time"
    "net"
)

type EchoWorker struct {
    Id          int
    TcpListener *net.TCPListener
    UdpConn     *net.UDPConn
    Sync        chan int
    IsTcp       bool
}

func NewEchoWorker(id int, listener interface{}, sync chan int) *EchoWorker {
    var worker  *EchoWorker

    if tcpListener, ok := listener.(*net.TCPListener); ok {
        worker  = &EchoWorker{
            Id:             id,
            TcpListener:    tcpListener,
            UdpConn:        nil,
            Sync:           sync,
            IsTcp:          true,
        }
    } else if udpConn, ok := listener.(*net.UDPConn); ok {
        worker  = &EchoWorker{
            Id:             id,
            TcpListener:    nil,
            UdpConn:        udpConn,
            Sync:           sync,
            IsTcp:          false,
        }
    }

    // currect the worker id by channel
    worker.Id, _ = <-worker.Sync
    return worker
}

func (w *EchoWorker) Run() {
    if w.IsTcp {
        w.RunTcp()
    } else {
        w.RunUdp()
    }
}

func (w *EchoWorker) RunTcp() {
    for {
        conn, err := w.TcpListener.Accept()
        if err != nil {
            continue
        }

        go w.handleClient(conn)
    }
}

func (w *EchoWorker) RunUdp() {
    var buf []byte = make([]byte, 512)

    for {
        time.Sleep(100 * time.Millisecond)
        _, address, _ := w.UdpConn.ReadFromUDP(buf)

        // 发送数据
        w.UdpConn.WriteToUDP(buf, address)
    }
}

func (w *EchoWorker) handleClient(client net.Conn) {
    defer client.Close()

    duration, _ := time.ParseDuration("500ms")
    client.SetDeadline(time.Now().Add(duration))

    var buf [512]byte
    n, _ := client.Read(buf[0:])

    client.Write(buf[0:n])
}

