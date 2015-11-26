// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package leveldb

import (
    "testing"
)

var (
    testdata    = []map[string]string{
        {
            "table1": "b1",
            "table2": "b2",
            "table3": "b3",
        },
        {
            "group1": "g1",
            "group2": "g2",
            "group3": "g3",
        },
    }

    dbfile  = "test.db"
)

func init() {
    DBRoot  = "/tmp"
}

// to test insert data by level Batch.
func Test_WriteLevelBatch(t *testing.T) {
    wb  := GetLevelBatch()
    wb.Put([]byte("tony"), []byte(EncodeData(testdata)))
    wb.Put([]byte("tony1"), []byte(EncodeData(testdata)))
    wb.Put([]byte("tony2"), []byte(EncodeData(testdata)))

    err := WriteBatch(wb, dbfile)
    if err != nil {
        t.Error(err)
    }
}

func Test_WriteLevelDB(t *testing.T) {

    err := WriteDB("tony", testdata, dbfile)
    WriteDB("tony1", testdata, dbfile)
    WriteDB("tony2", testdata, dbfile)
    if err != nil {
        t.Error(err)
    }
}

func Test_DeleteData(t *testing.T) {
    return
    wb  := GetBatch()
    wb.Delete([]byte("tony"))
    wb.Delete([]byte("tony1"))
    wb.Delete([]byte("tony2"))

    err := WriteBatch(wb, dbfile)
    if err != nil {
        t.Error(err)
    }
}

func Test_ReadLevelDBByte(t *testing.T) {
    _, err   := ReadDBByte(dbfile, nil, 0)
    if err != nil {
        t.Error(err)
    }
}

// to test get all data by limit
func Test_ReadLevelDBByLimit(t *testing.T) {
    _, err   := ReadDB(dbfile, "main", 0)
    if err != nil {
        t.Error(err)
    }
}

// to test get datat by condition
func Test_ReadLevelDBByCondition(t *testing.T) {
    _, err   := ReadDB(dbfile, "main", 0)
    if err != nil {
        t.Error(err)
    }
}

