// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package host

import (
    "fmt"
    "time"
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
	    t1 := time.Now().UnixNano()
        a    := Ping(ip, delay)
	    t2 := time.Now().UnixNano()
        if !a {
            fmt.Printf("Ths ip %s can not be araived!\n", ip)
            continue
        }

        fmt.Printf("Result: ip:%s, used time:%f ms\n", ip, float64(t2-t1) / 1000000)
    }
}
