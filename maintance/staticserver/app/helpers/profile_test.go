// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// profile相关帮助工具测试.

package helpers

import (
    "testing"
    "fmt"
)

func Test_GetBirthdayByAge(t *testing.T) {
    a   := GetBirthdayByAge(18)
    b   := a.Year()
    fmt.Printf("%#v, %d\n", a, b)
}

