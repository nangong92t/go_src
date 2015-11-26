package modules

import (
    "git.masontest.com/branches/goserver/plugins"
	rpc "git.masontest.com/branches/goserver/client"
    //"git.masontest.com/branches/goserver/apps/hoststatus/app/models"
)

var (
    logger *plugins.ServerLog
)

func init() {
    logger  = plugins.NewServerLog("tony")
}

func Run() {
    configs  := []map[string]string {{
        "Uri":          "10.161.185.137:8999",
        "User":         "tony",
        "Password":     "123",
        "Signature":    "ab1f8e61026a7456289c550cb0cf77cda44302b4",
    },
    {
        "Uri":          "127.0.0.1:8999",
        "User":         "tony",
        "Password":     "123",
        "Signature":    "ab1f8e61026a7456289c550cb0cf77cda44302b4",
    },
    }

    for i, config := range configs {
        conn    := rpc.NewRpcClient("udp", config)

        data    := conn.Call("host", "status", nil)
        logger.Add("server %d: %s", i, data)
    }
}
