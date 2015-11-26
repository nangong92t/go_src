package workers

import (
    "bufio"
    "fmt"
    "encoding/json"
//    "compress/gzip"
//    "bytes"
    "../protocols"
    "../plugins"
)

type TextWorker struct {
    Proto   *protocols.TextProtocol
    Log     *plugins.ServerLog
}

func NewTextWorker(logger *plugins.ServerLog) *TextWorker {
    return &TextWorker{
        Proto:  protocols.NewTextProtocol(logger),
        Log:    logger,
    }
}

func (w *TextWorker) Parse(reader *bufio.Reader) (map[string]string, error) {
    return  w.Proto.Input(reader)
}

func (w *TextWorker) Decode(data []byte) {

}

func (w *TextWorker) Encode(resp *protocols.Response) []byte {
    if resp.Mesg != "" {
        resp.Code   = 500
    }

    data, _ := json.Marshal(resp)

    /*
    var b bytes.Buffer
    gw   := gzip.NewWriter(&b)
    gw.Write(data)
    data    = b.Bytes()
    gw.Close()*/

    s   := fmt.Sprintf("%d\n%s\n", len(string(data)), string(data))

    return []byte(s)
}

