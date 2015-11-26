package protocols

type Request struct {
    Data        string          `json:"data"`
    Signature   string          `json:"signature"`
}

type RequestData struct {
    Version     string          `json:"version"`
    User        string          `json:"user"`
    Password    string          `json:"password"`
    Timestamp   float64         `json:"timestamp"`
    Class       string          `json:"class"`
    Method      string          `json:"method"`
    Params      []interface{}   `json:"params"`
}

type Response struct {
    Code        int             `json:"code"`
    Mesg        string          `json:"mesg"`
    Expend      float64         `json:"expend"`
    Data        interface{}     `json:"data"`
}

type Exception struct {
    Class       string          `json:"class"`
    Message     string          `json:"message"`
    Code        string          `json:"code"`
    File        string          `json:"file"`
    Line        string          `json:"line"`
    TraceAsString    string     `json:"traceasstring"`
}

type RemoteConfig struct {
    Version     string          `json:"version"`
    User        string          `json:"user"`
    Password    string          `json:"password"`
    Timestamp   int64           `json:"timestamp"`
    Class       string          `json:"class"`
    Method      string          `json:"method"`
    Params      []interface{}   `json:"params"`
    Signature   string          `json:"-"`
}

