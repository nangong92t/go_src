package protocols

import (
    "bufio"
    "errors"
    "strconv"
    "strings"
    "../plugins"
)

type TextProtocol struct {
    Log     *plugins.ServerLog
}

func NewTextProtocol(logger *plugins.ServerLog) *TextProtocol {
    return &TextProtocol{
        Log: logger,
    }
}

func (t *TextProtocol) Input(reader *bufio.Reader) (map[string]string, error) {
    var cmdlen, datalen int
    var cmd, data string

    result  := map[string]string{}

    errMsg  := errors.New("protocol error!")

    for i:=1; i<5; i++ {
        line, err := reader.ReadBytes('\n')
        if err != nil { // EOF, or worse
            return nil, errMsg
        }

        lineStr     := strings.Replace(
            strings.Replace(string(line), "\r\n", "", -1),
            "\n", "", -1,
        )

        switch i {
            case 1:
                cmdlen  = t.getLength(lineStr)
                // t.Log.Add("protocol 1: ", cmdlen, lineStr)
                if cmdlen == 0 { return nil, errMsg }
            case 2:
                cmd     = t.getData(lineStr, cmdlen)
                // t.Log.Add("protocol 2: %s", cmd, lineStr)
                if cmd == "" { return nil, errMsg }
            case 3:
                datalen = t.getLength(lineStr)
                // t.Log.Add("protocol 3: %s", datalen, lineStr)
                if datalen == 0 { return nil, errMsg }
            case 4:
                data    = t.getData(lineStr, datalen)
                // t.Log.Add("protocol 4: %s", data, lineStr)
                if data == "" { return nil, errMsg }
        }
    }

    result["cmd"]   = cmd
    result["data"]  = data

    return result, nil
}

func (t *TextProtocol) getLength(data string) int {
    if length, err := strconv.Atoi(data); err == nil {
        return length
    } else {
        return 0
    }
}

func (t *TextProtocol) getData(data string, length int) string {
    if len(data) == length {
        return data
    } else {
        return ""
    }
}


