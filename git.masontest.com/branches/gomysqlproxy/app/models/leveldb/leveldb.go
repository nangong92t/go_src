// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package leveldb

import (
    "fmt"
    "time"
    "errors"
    "encoding/json"
    "encoding/hex"
    "crypto/md5"
    "math/rand"
    "github.com/jmhodges/levigo"
)

var (
    DBRoot  string
    DBConns = map[string]*levigo.DB{}

    DbName  = map[string]string{
        "MysqlProxy":       "main.db",
        "MysqlTable":       "table.db",
        "MysqlShardTable":  "shardtable.db",
        "MysqlShardDB":     "sharddb.db",
    }
)

type kv struct {
    K []byte
    V interface{}
}

type kvs struct {
    items map[int]kv
}

func (p *kv) PutKV(k []byte, v interface{}) {
    p.K = k
    p.V = v
}

func (items *kvs) PutKVs() {
    fmt.Println(items)
}

func (p *kv) GetKV() (key []byte, value interface{}) {
    key = p.K
    value = p.V
    return
}

func openDB(dbfile string) (*levigo.DB, error) {
    if conn, isOk := DBConns[dbfile]; isOk {
        return conn, nil
    }
    opts    := levigo.NewOptions()
    opts.SetCache(levigo.NewLRUCache(3<<20))
    opts.SetCreateIfMissing(true)
    db, err := levigo.Open(dbfile, opts)

    if err != nil { return nil, err }

    DBConns[dbfile] = db

    return db, nil
}

func getDbName(as string) (string, error) {
    if name, ok := DbName[as]; ok {
        return name, nil
    }

    return "", errors.New("Sorry, no this database alias in leveldb")
}

func GetLevelBatch() *levigo.WriteBatch {
    return levigo.NewWriteBatch()
}

func EncodeData(data interface{}) string {
    res, _ := json.Marshal(data)

    return string(res)
}

func DecodeData(data string) map[string]interface{} {
    var res map[string]interface{}
    json.Unmarshal([]byte(data), &res)

    return res
}

func UpdateDB(k string, v interface{}, dbName string) error {
    dbfile, err := getDbName(dbName)
    if err != nil { return err }

    dbfile  = DBRoot + "/" + dbfile

    db, err := openDB(dbfile)
    if err != nil { return err }

    data, err := json.Marshal(v)
    if err != nil { return err }

    wo := levigo.NewWriteOptions()
    err = db.Delete(wo, []byte(k))
    if err != nil { return err }

    err = db.Put(wo, []byte(k), data)
    if err != nil { return err }

    //db.Close()
    return nil
}

func WriteDB(k string, v interface{}, dbName string) error {
    dbfile, err := getDbName(dbName)
    if err != nil { return err }

    dbfile  = DBRoot + "/" + dbfile

    db, err := openDB(dbfile)
    if err != nil { return err }

    data, err := json.Marshal(v)
    if err != nil { return err }

    wo := levigo.NewWriteOptions()
    err = db.Put(wo, []byte(k), data)
    if err != nil { return err }

    //db.Close()
    return nil
}

func WriteBatch(wb *levigo.WriteBatch, dbName string) error {
    dbfile, err := getDbName(dbName)
    if err != nil { return err }

    dbfile  = DBRoot + "/" + dbfile

    db, err := openDB(dbfile)
    if err != nil { return err }

    wo := levigo.NewWriteOptions()
    err = db.Write(wo, wb)

    //db.Close()
    return err
}

// Get the data count from leveldb
//
// @param string dbName The levelDb database name.
// @param func   filter The condition filter function.
//
func ReadDBCount(dbName, condition string) (int, error) {
    dbfile, err := getDbName(dbName)
    if err != nil { return 0, err }

    dbfile  = DBRoot + "/" + dbfile

    db, err := openDB(dbfile)
    if err != nil { return 0, err }

    ro := levigo.NewReadOptions()
    ro.SetFillCache(false)
    it := db.NewIterator(ro)
    defer it.Close()

    foundCnt    := 0
    for it.Seek([]byte(condition)); it.Valid(); it.Next() {
        foundCnt++
    }

    if err := it.GetError(); err != nil {
        return 0, err
    }

    return foundCnt, nil
}


// Get the data from leveldb
//
// @param string dbName The levelDb database name.
// @param func   filter The condition filter function.
// @parrm int    limit  The limit result items number, if limit==0, means no any limit. 
//
func ReadDB(dbName, condition string, limit int) (items []map[string]interface{}, err error) {
    dbfile, err := getDbName(dbName)
    if err != nil { return nil, err }

    dbfile  = DBRoot + "/" + dbfile

    db, err := openDB(dbfile)
    if err != nil { return nil, err }

    ro := levigo.NewReadOptions()
    ro.SetFillCache(false)
    it := db.NewIterator(ro)
    defer it.Close()

    foundCnt    := 0
    for it.Seek([]byte(condition)); it.Valid(); it.Next() {
        var data interface{}
        k := it.Key()

        foundCnt++
        v := it.Value()
        json.Unmarshal(v, &data)

        items   = append(items, map[string]interface{}{
            string(k): data,
        })

        if limit > 0 && limit == foundCnt { break }
    }

    if err := it.GetError(); err != nil {
        return nil, err
    }

    return items, nil
}

// Get the data from leveldb
//
// @param string dbName The levelDb database file path.
// @param func   filter The condition filter function.
// @parrm int    limit  The limit result items number, if limit==0, means no any limit. 
//
func ReadDBByte(dbName string, filter func([]byte) bool, limit int) (map[int]kv, error) {
    dbfile, err := getDbName(dbName)
    if err != nil { return nil, err }

    dbfile  = DBRoot + "/" + dbfile

    db, err := openDB(dbfile)
    if err != nil { return nil, err }

    ro := levigo.NewReadOptions()
    ro.SetFillCache(false)
    it := db.NewIterator(ro)
    defer it.Close()

    foundCnt    := 0

    //items := map[int64]kv{}
    item := new(kv)
    items := map[int]kv{}

    for it.Seek([]byte("key")); it.Valid(); it.Next() {
        var data interface{}
        k := it.Key()
        if filter != nil && !filter(k) { continue }

        v := it.Value()
        json.Unmarshal(v, &data)
        item.PutKV(k, data)
        items[foundCnt] = *item

        if limit > 0 && limit == (foundCnt+1) { break }
        foundCnt++
    }

    if err := it.GetError(); err != nil {
        return nil, err
    }

    return items, nil
}

// To build the md5 primary key.
//
// @param string key    To build the key according to this key.
// @param bool   isRand If it's true the primary key will change for each call.
//
func BuildPrimaryKey(key string, isRand bool) string {
    h       := md5.New()
    if isRand {
        curT    := time.Now().UnixNano()
        ra      := rand.New(rand.NewSource(curT))
        key     = key + fmt.Sprintf("%s", curT) + fmt.Sprintf("%s", ra.Intn(1000))
    }
    h.Write([]byte(key))
    return fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
}

func GetCurTime() int64 {
    return time.Now().Unix()
}

// change interface to id string
func RestorePrimaryId(iids []interface{}) []string {
    ids     := []string{}
    for _, id := range iids { ids = append(ids, id.(string)) }
    return ids
}


