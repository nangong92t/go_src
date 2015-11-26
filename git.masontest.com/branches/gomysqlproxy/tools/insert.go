package main

import (
    "time"
    "fmt"
    "runtime"
    "strconv"

    rpc "git.masontest.com/branches/goserver/client"
)

type result struct {
    Id int
    Tablename string
    Runt float64
    Data interface{}
}

var (
    conn *rpc.RpcClient
)

func init() {
    config  := map[string]string {
        "Uri":          "121.199.44.61:3307",
        "User":         "tony",
        "Password":     "123",
        "Signature":    "ab1f8e61026a7456289c550cb0cf77cda44302b4",
    }
    conn    = rpc.NewRpcClient("tcp", config)
}

func insertTable() {
    t1  := time.Now().UnixNano()
    data, _    := conn.Call("proxy", "exec", []interface{}{[]string{"insert into maybe values('', ?)", "TonyXu"}})
    t2  := time.Now().UnixNano()
    fmt.Printf("data: %#v\n time: %4.3f s\n", data, float64(t2-t1)/1000000000)
}

func createTable(i int) *result {
    tableName   := "test_table_" + strconv.Itoa(i)

fmt.Printf("table name: %s, id: %d\n", tableName, i)
    t1  := time.Now().UnixNano()
    curData, err    := conn.Call("proxy", "exec", []interface{}{[]string{"create table "+tableName+"( id int(11) not null auto_increment primary key, name varchar(30) not null)"}})
    if err != nil { fmt.Printf("table %s has problem: %s\n", tableName, err); return nil }
    t2  := time.Now().UnixNano()

    return &result{
        Id:     i,
        Tablename: tableName,
        Runt:   float64(t2-t1)/1000000000,
        Data:   curData,
    }
}

func alterTable(i int) *result {
    tableName   := "test_table_" + strconv.Itoa(i)

fmt.Printf("table name: %s, id: %d\n", tableName, i)
    t1  := time.Now().UnixNano()
    curData, err    := conn.Call("proxy", "exec", []interface{}{[]string{"alter table "+tableName+" change name name2 varchar(40) not null"}})
    if err != nil { fmt.Printf("table %s has problem: %s\n", tableName, err); return nil }
    t2  := time.Now().UnixNano()

    return &result{
        Id:     i,
        Tablename: tableName,
        Runt:   float64(t2-t1)/1000000000,
        Data:   curData,
    }
}

func dropTable(i int) *result {
    tableName   := "test_table_" + strconv.Itoa(i)

fmt.Printf("table name: %s, id: %d\n", tableName, i)
    t1  := time.Now().UnixNano()
    curData, err    := conn.Call("proxy", "exec", []interface{}{[]string{"drop table "+tableName}})
    if err != nil { fmt.Printf("table %s has problem: %s\n", tableName, err); return nil }
    t2  := time.Now().UnixNano()

    return &result{
        Id:     i,
        Tablename: tableName,
        Runt:   float64(t2-t1)/1000000000,
        Data:   curData,
    }
}

func insertData(i int) *result {
    tableName   := "test_table_i4"
    curName     := "tonyTest_" + strconv.Itoa(i)

fmt.Printf("data: %s, id: %d\n", curName, i)
    t1  := time.Now().UnixNano()
    curData, err    := conn.Call("proxy", "exec", []interface{}{[]string{"insert into "+tableName + " values('', ?)", curName}})
    if err != nil { fmt.Printf("table %s has problem: %s\n", tableName, err); return nil }
    t2  := time.Now().UnixNano()

    return &result{
        Id:     i,
        Tablename: tableName,
        Runt:   float64(t2-t1)/1000000000,
        Data:   curData,
    }
}

func getStatus() {
    curData, _    := conn.Call("proxy", "getstatus", nil)
    fmt.Printf("current status: %s", curData)
}

func testCreate() {
    runCh       := make(chan *result, runtime.NumCPU())
    runTotal    := 5    //1000
    maxRunt     := float64(0)
    minRunt     := float64(999.0)
    avgRunt     := float64(0)
    allRunt     := float64(0)

    for i:=0; i<runTotal; i++ {
        go func(curI int) {
            runCh <- insertData(curI) //createTable(curI)
        }(i)
    }

    for i:=0; i<runTotal; i++ {
        res := <-runCh
        if res == nil { continue }

        fmt.Printf("id: %d, table name: %s, runt: %4.3f, data: %#v\n", res.Id, res.Tablename, res.Runt, res.Data)
        if res.Runt > maxRunt { maxRunt = res.Runt }
        if res.Runt < minRunt { minRunt = res.Runt }
        allRunt += res.Runt
    }

    avgRunt = allRunt / float64(runTotal)

    fmt.Printf("run times: %d, min: %4.3f, max: %4.3f, avg: %4.3f\n", runTotal, minRunt, maxRunt, avgRunt)
}

func main() {
    testCreate()
    // getStatus()
}

