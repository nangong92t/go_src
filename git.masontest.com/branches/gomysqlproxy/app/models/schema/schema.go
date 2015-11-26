// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

package schema

// Yes, this sucks. It's a tiny tiny package that needs to be on its own
// It contains a data structure that's shared between sqlparser & tabletserver

import (
    "fmt"
	"strings"

    "git.masontest.com/branches/gomysqlproxy/app/models/sqltypes"
)

// Column categories
const (
	CAT_OTHER = iota
	CAT_NUMBER
	CAT_VARBINARY
)

// Cache types
const (
	CACHE_NONE = 0
	CACHE_RW   = 1
	CACHE_W    = 2
)

type TableColumn struct {
	Name     string
	Category int
	IsAuto   bool
	Default  sqltypes.Value
    Key      string
    Extra    string
    Null     bool
}

func (ta *MysqlTable) AddColumn(name string, columnType string, defval sqltypes.Value, extra string) {
	index := len(ta.Columns)
	ta.Columns = append(ta.Columns, TableColumn{Name: name})
	if strings.Contains(columnType, "int") {
		ta.Columns[index].Category = CAT_NUMBER
	} else if strings.HasPrefix(columnType, "varbinary") {
		ta.Columns[index].Category = CAT_VARBINARY
	} else {
		ta.Columns[index].Category = CAT_OTHER
	}
	if extra == "auto_increment" {
		ta.Columns[index].IsAuto = true
		// Ignore default value, if any
		return
	}
	if defval.IsNull() {
		return
	}
	if ta.Columns[index].Category == CAT_NUMBER {
		ta.Columns[index].Default = sqltypes.MakeNumeric(defval.Raw())
	} else {
		ta.Columns[index].Default = sqltypes.MakeString(defval.Raw())
	}
}

func (tbl *MysqlTable) RemoveColumn(name string) {
    newColumns  := []TableColumn{}
    for _, col  := range tbl.Columns {
        if col.Name != name {
            newColumns  = append(newColumns, col)
        }
    }

    tbl.Columns = newColumns
}

// about the params detail infomation, please to check  sqlparser/functions_test.go 
func (tbl *MysqlTable) ChangeColumn(name string, params map[string]string) {
    var curCol *TableColumn
    var foundIdx    int

    for i, col  := range tbl.Columns {
        if col.Name == name {
            curCol  = &col
            foundIdx    = i
            break
        }
    }

    if curCol == nil { return }

    if _, isOk  := params["newName"]; isOk {
        curCol.Name = params["newName"]
    }

    if _, isOk  := params["auto_increment"]; isOk {
        curCol.IsAuto   = true
    }

    if _, isOk  := params["null"]; isOk {
        if _, isNot := params["not"]; isNot {
            curCol.Null = false
        } else {
            curCol.Null = true
        }
    }

    tbl.Columns[foundIdx]   = *curCol
}

// to restore the table column to memery.
//
func (ta *MysqlTable) RestoreColumnsByDB() error {
    shardTabOrderId := 0
    db, err := ta.GetSlaveShardDBConn(shardTabOrderId)
    if err != nil { return err }

    shardTblN   := ta.Shards[shardTabOrderId].Name
    // var rows *sql.Rows
    rows, err := db.Query(fmt.Sprintf("desc `%s`", shardTblN))
    if err != nil { return err }

    if rows != nil {
        for rows.Next() {
            var colName string
            var colType string
            var colNull string
            var colKey []byte
            var colDefault []byte
            var colExtra string

            err     = rows.Scan(&colName, &colType, &colNull, &colKey, &colDefault, &colExtra)
            if err != nil { return err }

            // if colNull=="" || string(colKey)=="" {}
            ta.AddColumn(colName, colType, sqltypes.MakeString(colDefault), colExtra)
        }

        rows.Close()
    }

    return nil
}

func (ta *MysqlTable) FindColumn(name string) int {
	for i, col := range ta.Columns {
		if col.Name == name {
			return i
		}
	}
	return -1
}

func (ta *MysqlTable) GetPKColumn(index int) *TableColumn {
	return &ta.Columns[ta.PKColumns[index]]
}

func (ta *MysqlTable) GetAutoColumn() *TableColumn {
    for _, col  := range ta.Columns {
        if col.IsAuto {
            return &col
        }
    }
    return nil
}

func (ta *MysqlTable) AddIndex(name string) (index *Index) {
	index = NewIndex(name)
	ta.Indexes = append(ta.Indexes, index)
	return index
}

type Index struct {
	Name        string
	Columns     []string
	Cardinality []uint64
	DataColumns []string
}

func NewIndex(name string) *Index {
	return &Index{name, make([]string, 0, 8), make([]uint64, 0, 8), nil}
}

func (idx *Index) AddColumn(name string, cardinality uint64) {
	idx.Columns = append(idx.Columns, name)
	if cardinality == 0 {
		cardinality = uint64(len(idx.Cardinality) + 1)
	}
	idx.Cardinality = append(idx.Cardinality, cardinality)
}

func (idx *Index) FindColumn(name string) int {
	for i, colName := range idx.Columns {
		if name == colName {
			return i
		}
	}
	return -1
}

func (idx *Index) FindDataColumn(name string) int {
	for i, colName := range idx.DataColumns {
		if name == colName {
			return i
		}
	}
	return -1
}
