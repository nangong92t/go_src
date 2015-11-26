package client

import (
    "testing"
    "fmt"

    rpc "git.masontest.com/branches/goserver/client"
)

func Test_RpcClient(t *testing.T) {
    config  := map[string]string {
        "Uri":          "127.0.0.1:9901",
        "User":         "tony",
        "Password":     "123",
        "Signature":    "123",
    }
    conn    := rpc.NewRpcClient("udp", config)

    fmt.Printf("%#v", conn)

}

func Test_ErrorRpcClient(t *testing.T) {
    defer func() {
        if err := recover(); err != nil {
            if fmt.Sprintf("%s", err) != "Please config the value for key: Signature" {
                t.Error("error type wrong!")
            }
        }
    }()
    config  := map[string]string {
        "Uri":          "127.0.0.1:9901",
        "User":         "tony",
        "Password":     "123",
    }
    conn    := rpc.NewRpcClient("udp", config)

    fmt.Printf("%#v", conn)
}

func Test_RemoteCall(t *testing.T) {
    config  := map[string]string {
        "Uri":          "127.0.0.1:9901",
        "User":         "tony",
        "Password":     "123",
        "Signature":    "ab1f8e61026a7456289c550cb0cf77cda44302b4",
    }
    conn    := rpc.NewRpcClient("udp", config)

    data, _    := conn.Call("host", "ping", nil)
    panic(fmt.Sprintf("%#v", data))
    isOk    := false

    switch data.Data.(type) {
        case string:
            if data.Data.(string) == "Pong" {
                isOk    = true
            }
    }

    if !isOk {
        t.Error("test call failed!")
    }
}

func Test_RemoteCallStatus(t *testing.T) {
    config  := map[string]string {
        "Uri":          "127.0.0.1:9901",
        "User":         "tony",
        "Password":     "123",
        "Signature":    "ab1f8e61026a7456289c550cb0cf77cda44302b4",
    }
    conn    := rpc.NewRpcClient("udp", config)

    data, _    := conn.Call("host", "status", nil)
    panic(fmt.Sprintf("%#v", data))
}
