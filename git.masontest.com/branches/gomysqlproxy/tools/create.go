package main

import (
    "time"
    "fmt"

    rpc "git.masontest.com/branches/goserver/client"
)

func main() {
    t1  := time.Now().UnixNano()
    config  := map[string]string {
        "Uri":          "121.199.44.61:3307",
        "User":         "tony",
        "Password":     "123",
        "Signature":    "ab1f8e61026a7456289c550cb0cf77cda44302b4",
    }
    conn    := rpc.NewRpcClient("tcp", config)

    data, _    := conn.Call("proxy", "exec", []interface{}{[]string{"CREATE  TABLE IF NOT EXISTS `album` (`album_id` INT NOT NULL AUTO_INCREMENT primary key,    `cover_pictureid` int not null default 0,    `album_name` char(60) NOT NULL default '',    `album_desc` char(200) NOT NULL default '',    `created_time` int not null,    `creator` int(11) not null comment '创建者id',    `status` tinyint(1) not null default 0,    `picture_total` int not null default 0,    `is_allow_comment` tinyint(1) not null default 0,  UNIQUE INDEX `name_UNIQUE` (`album_name` ASC) )"}})
    t2  := time.Now().UnixNano()
    fmt.Printf("data: %#v\n time: %4.3f s\n", data, float64(t2-t1)/1000000000)
}

