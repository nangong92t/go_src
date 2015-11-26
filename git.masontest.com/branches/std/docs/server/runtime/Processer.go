package runtime

import (
    "reflect"
    "../protocols"
)

func Processer(appName string, req *protocols.RequestData) (interface{}, error) {
    appName := "Process" + appName
    return reflect.TypeOf((*websocket.Conni)(nil))
}
