// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package utils

import (
    "fmt"
    "errors"
)

func EToI64(old string) (uint64, error) {
    var newI float64
    n, err := fmt.Sscanf(old, "%e", &newI)
    if err != nil {
        return uint64(0), err
    } else if 1 != n {
        return uint64(0), errors.New("translate error!")
    }

    return uint64(newI), nil
}


