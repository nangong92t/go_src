package protocols

type RequestData struct {
    Version     string          `json:"version"`
    Class       string          `json:"class"`
    Method      string          `json:"method"`
    Params      []interface{}   `json:"params"`
    Signature   string          `json:"signature"`
}

type Response struct {
    Code        int             `json:"code"`
    Mesg        string          `json:"mesg"`
    Data        interface{}     `json:"data"`
}

type Exception struct {
    Class       string          `json:"class"`
    Message     string          `json:"message"`
    Code        string          `json:"code"`
    File        string          `json:"file"`
    Line        string          `json:"line"`
    TraceAsString    string  `json:"traceasstring"`
}

