// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// To defined global json render.

package results

import (
)

type RenderErrorResult struct {
    Error       string  `json:"error"`
    ErrorCode   int     `json:"errorcode"`
}

type RenderApiResult struct {
    Result        interface{} `json:"result"`
}

type RenderApiErrResult struct {
    Errors      []*RenderErrorResult `json:"errors"`
}
