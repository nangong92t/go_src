// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package models

import (
    "fmt"
    "testing"
)

func Test_Ping(t *testing.T) {
    delay   := 1000
    hosts   := []string{
        "112.124.109.11",
        "114.215.184.73",
        "182.92.130.30",
        "123.57.37.140",
        "112.124.114.237",
        "182.92.150.51",
    }

    for _, ip := range hosts {
        t, e  := Ping(ip, delay)
        if e != nil {
            fmt.Printf("Ths ip %s can not be araived!, error: %s\n", ip, e)
            continue
        }

        fmt.Printf("Result: ip:%s, used time:%f ms\n", ip, float64(t) / 1000000)
    }
}
