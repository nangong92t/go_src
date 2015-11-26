package main

import (
    "fmt"
    "strconv"
    //"testing"
    "log"
    "time"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

func ConnToDB() (*sql.DB, error) {
    connStrFmt  := "%s:%s@tcp([%s]:%s)/%s?charset=utf8"
    // if isWrite { connStrFmt += "&autocommit=true" }
    connStr := fmt.Sprintf(connStrFmt, "root", "shaluo", "121.199.44.61", "3306", "shard2")
    db, err := sql.Open("mysql", connStr)

    if err != nil {
        return nil, err
    }

    return db, nil
}

func mysqlCreate() {
    db, err := ConnToDB()
    if err != nil { log.Fatalf("%s", err) }

    stmt, err   := db.Prepare("insert into maybe(name) values(?)")
    if err != nil { log.Fatalf("%s", err) }

    _, err    = stmt.Exec("Tonyhaha1")
    if err != nil { log.Fatalf("%s", err) }

    stmt.Close()
    db.Close()
}

func dropTable() {
    db, err := ConnToDB()
    if err != nil { log.Fatalf("%s", err) }

    for i:=0; i<500; i++ {
        stmt, err   := db.Prepare("drop table test1_table_"+ strconv.Itoa(i)+ "_shard1")
        if err != nil { log.Fatalf("%s", err) }

        _, err    = stmt.Exec()
        if err != nil { log.Fatalf("%s", err) }

        stmt.Close()
    }
    db.Close()
}

func main() {
    t1  := time.Now().UnixNano()
    dropTable()
    t2  := time.Now().UnixNano()

    fmt.Printf("%4.3f s\n", float64(t2-t1)/1000000000)
}

