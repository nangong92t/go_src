package protocols

import (
    "bufio"
    "errors"
    "time"
    "strings"
    "crypto/md5"
    "strconv"
    "fmt"
    "encoding/json"
    "git.masontest.com/branches/goserver/plugins"
)

var (
    // curent version.
    Version = "1.0"

    // current support cmds
    Cmds    = map[string]int{
        "RPC": 1,
        "RPC:GZ": 1,
    }
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
            t.Log.Add("protocol 0: %s", err)
            return nil, err
        }

        lineStr     := strings.Replace(
            strings.Replace(string(line), "\r\n", "", -1),
            "\n", "", -1,
        )

        switch i {
            case 1:
                cmdlen  = t.getLength(lineStr)
                if cmdlen == 0 {
                    t.Log.Add("protocol error 1: ", cmdlen, lineStr)
                    return nil, errMsg
                }
            case 2:
                cmd     = t.getData(lineStr, cmdlen)
                if cmd == "" {
                    t.Log.Add("protocol error 2: %s", cmd, lineStr)
                    return nil, errMsg
                }
            case 3:
                datalen = t.getLength(lineStr)
                if datalen == 0 {
                    t.Log.Add("protocol error 3: %s", datalen, lineStr)
                    return nil, errMsg
                }
            case 4:
                data    = t.getData(lineStr, datalen)
                if data == "" {
                    t.Log.Add("protocol error 4: %s", data, lineStr)
                    return nil, errMsg
                }
        }
    }

    result["cmd"]   = cmd
    result["data"]  = data

    return result, nil
}

func (t *TextProtocol) ClientInput(dataSource []byte) (*Response, error) {
    var datalen int
    var data string

    resp    := &Response{}

    errMsg  := errors.New("protocol error in client!")

    datas   := strings.Split(string(dataSource), "\n")
    if len(datas) != 3 {
        t.Log.Add("protocol error in client!")
        return nil, errMsg
    }

    for i:=0; i<2; i++ {
        lineStr     := strings.Replace(
            strings.Replace(datas[i], "\r\n", "", -1),
            "\n", "", -1,
        )

        switch i {
            case 0:
                datalen  = t.getLength(lineStr)
                if datalen == 0 {
                    t.Log.Add("protocol error 1 in client: %s", datalen, lineStr)
                    return nil, errMsg
                }
            case 1:
                data    = t.getData(lineStr, datalen)
                if data == "" {
                    t.Log.Add("protocol error 2 in client: %s", data, lineStr)
                    return nil, errMsg
                }
        }
    }

    err     := json.Unmarshal([]byte(data), &resp)
    if err != nil { t.Log.Add("protocol error 2 in client: %s", data); return nil, errMsg }

    return resp, nil
}

func (t *TextProtocol) GetActiveSignature(secretKey, data string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(data + "&" + secretKey)))
}

func (t *TextProtocol) getLength(data string) int {
    if length, err := strconv.Atoi(data); err == nil {
        return length
    } else {
        return 0
    }
}

func (t *TextProtocol) Decode(data []byte) ([]byte, error) {
    data, err := plugins.GzipDecode(data)
    return data, err
}

func (t *TextProtocol) Encode(resp *Response, cmd string) []byte {
    if resp.Mesg != "" {
        resp.Code   = 500
    }

    data, _ := json.Marshal(resp)

    if cmd == "RPC:GZ" {
        data, _ = plugins.GzipEncode(data)
    }
    s   := fmt.Sprintf("%d\n%s\n", len(data), data)

    return []byte(s)
}


func (t *TextProtocol) getData(data string, length int) string {
    if len(data) == length {
        return data
    } else {
        return ""
    }
}

// to package the text protocol.
// @param cmd, at parent can support : RPC, RPC:GZ
func (t *TextProtocol) Output(cmd string, innerData *RemoteConfig) ([]byte, error) {
    if _, ok := Cmds[cmd]; !ok {
        return nil, errors.New("sorry, GoServer don't support this cmd")
    }

    innerData.Timestamp = time.Now().Unix()

    data, err   := json.Marshal(innerData)
    if err != nil { return nil, err }
    dataStr     := string(data)

    signature   := t.GetActiveSignature(innerData.Signature, dataStr)
    data, err   = json.Marshal(Request{
        Data:       dataStr,
        Signature:  signature,
    })
    if err != nil { return nil, err }

    dataStr     = string(data)
    outData     := []byte(
        strconv.Itoa(len(cmd))+"\n"+cmd+"\n"+strconv.Itoa(len(dataStr))+"\n"+dataStr+"\n",
    )

    return outData, err
}

// build a new Remote  config struct
func NewRemoteConfig(config map[string]string) (*RemoteConfig, error) {
    needFields  := []string{
        "User", "Password", "Signature",
    }

    for _, key := range needFields {
        if _, ok := config[key]; !ok {
            return nil, errors.New("Please config the value for key: " + key)
        }
    }

    t   := &RemoteConfig{
        Version:    Version,
        User:       config["User"],
        Password:   config["Password"],
        Timestamp:  int64(0),
        Class:      "",
        Method:     "",
        Params:     []interface{}{},
        Signature:  config["Signature"],
    }

    return t, nil
}

func (r *RemoteConfig) CleanParams() {
    r.Class     = ""
    r.Method    = ""
    r.Params    = []interface{}{}
}
