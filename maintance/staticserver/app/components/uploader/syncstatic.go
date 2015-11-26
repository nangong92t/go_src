// Copyright 2015, WePiao Teq Inc. All rights reserved.
// Author XiRongYi<xurongyi@wepiao.com>,
// Build on dev-0.0.1
// MIT Licensed

package uploader

import (
    "fmt"
    "bytes"
    "mime/multipart"
    "net"
    "net/http"
    "io"
    "os"
    "strconv"
    "strings"
    "github.com/revel/revel"
)

var StaticHosts = []string{}
var LocalHost   = ""
var ServerPort  = ""

// 等待同步的图片(绝对路径)队列
var WaitSyncFiles   = []string{}

// 同步失败数据模型.
type SyncFailed struct {
    Try     int
    Host    string
    Path    string
    Err     error
}

// 等待再次尝试上次的失败图片队列.
var WaitSyncFailedFiles = []SyncFailed{}

// 当前是否正在执行同步. 同时只能存在1次同步动作.
var IsSyncing = false

// 当前是否正在执行再次同步.
var IsTrySyncing    = false

func init() {
    // just test.
    // WaitSyncFiles   = append(WaitSyncFiles, "/home/tony/code/static/upload/9/6f2/845de/96f2c845ded8dfd842c5da24bda1bda5.jpg")

    LocalHost = GetCurIP()
}

func GetCurIP() string {
    conn, err := net.Dial("udp", "baidu.com:80")
    if err != nil {
        fmt.Println(err.Error())
        return ""
    }
    defer conn.Close()

    return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func GetParamString(param string, defaultValue string) string {
    revel.Config.SetSection( revel.RunMode )
    p, found := revel.Config.String(param)
    if !found {
        if defaultValue == "" {
            revel.ERROR.Fatal("Cound not find parameter: " + param)
        } else {
            return defaultValue
        }
    }

    return p
}

func GetAllStaticHosts() error {
    staticTotal, _ := strconv.Atoi(GetParamString("static.host.total", "0"))
    if staticTotal == 0 { return nil }

    ServerPort      = GetParamString("http.port", "9001")

    for i:=0; i<staticTotal; i++ {
        curHost := GetParamString("static." + strconv.Itoa(i) + ".host", "")
        if curHost != "" {
            StaticHosts = append(StaticHosts, curHost)
        }
    }
    return nil
}

func SyncNewToAllStaticHosts() error {
    if IsSyncing { return nil }
    IsSyncing   = true
    defer func() { IsSyncing = false }()

    syncFiles   := WaitSyncFiles
    if len(syncFiles) <= 0 { return nil }
    WaitSyncFiles   = []string{}

    runChan     := make(chan error)
    goTotal     := 0
    for _, curFile := range syncFiles {
        for _, curHost := range StaticHosts {
            if curHost == LocalHost { continue }
            go uploadFileToOther(curHost, curFile, runChan)
            goTotal++
        }
    }

    for i:=0; i<goTotal; i++ {
        curErr  := <-runChan
        if curErr != nil {
            fmt.Printf("upload error: %#v\n", curErr)
        }
    }

    IsSyncing   = false
    return nil
}

func TrySyncFailedStatic() error {
    if IsTrySyncing { return nil }
    IsTrySyncing    = true

    syncFailedFiles := WaitSyncFailedFiles
    if len(syncFailedFiles) <= 0 { return nil }
    WaitSyncFailedFiles = []SyncFailed{}

    IsTrySyncing    = false
    return nil
}

// 执行同步功能.
func uploadFileToOther(host, file string, errCh chan error) error {
    var syncUrl = "http://" + host + ":" + ServerPort + "/upload-sync"
    var b bytes.Buffer
    w   := multipart.NewWriter(&b)

    // read file.
    f, err  := os.Open(file)
    if err != nil {
        errCh <- err
        return err
    }

    fw, err := w.CreateFormFile("file", file)
    if err != nil {
        errCh <- err
        return err
    }

    if _, err = io.Copy(fw, f); err != nil {
        errCh <- err
        return err
    }
    w.Close()
    req, err    := http.NewRequest("POST", syncUrl, &b)
    if err != nil {
        errCh <- err
        return err
    }

    req.Header.Set("Content-Type", w.FormDataContentType())

    client  := &http.Client{}
    res, err    := client.Do(req)
    if err != nil {
        errCh <- err
        return err
    }

    // Check the response.
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
        errCh <- err
        return err
    }

    errCh <- nil
    return nil
}
