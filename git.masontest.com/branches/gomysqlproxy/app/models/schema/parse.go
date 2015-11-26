// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// parse the sql submodel.

package schema

import (
    "fmt"
    "strings"
    "errors"

    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

func (tbl *MysqlTable) ParseMergeInsertGlobalId(sqlSource []string, shardTable *MysqlShardTable) (string, []interface{}, error) {
    sql     := sqlSource[0]
    params  := []interface{}{}
    for _, one := range sqlSource[1:] {
        params  = append(params, one)
    }

    autoCol := tbl.GetAutoColumn()

    // no anyone auto increment column.
    if autoCol == nil { return sqlSource[0], params, nil }

    cols, vals, other, typ  := sqlparser.ParseInsertPart(sql, tbl.Name)

    autoIndex   := -1
    sqlTpl      := ""

    switch typ {
    case 1, 3:
        for i, one  := range cols {
            cols[i] = strings.TrimSpace(one)
            if cols[i] == autoCol.Name { autoIndex = i }
        }
        for i, one  := range vals {
            if autoIndex == i {
                vals[i] = tbl.GetGId()
            } else {
                vals[i] = strings.TrimSpace(one)
            }
        }
        if autoIndex == -1 {
            cols    = append(cols, autoCol.Name)
            vals    = append(vals, tbl.GetGId())
        }

        sqlTpl      = "insert into " + shardTable.Name + "(%s)"

        if typ == 1 {
            sqlTpl  += " select %s from %s"
            // TODO: need to parse the select other from table shard db.....

            sql     = fmt.Sprintf(sqlTpl, strings.Join(cols, ", "), strings.Join(vals, ", "), other)
        } else {
            sqlTpl  += " values(%s)"
            sql     = fmt.Sprintf(sqlTpl, strings.Join(cols, ", "), strings.Join(vals, ", "))
        }

    case 2, 4:
        for i, _ := range tbl.Columns {
            if tbl.Columns[i].Name == autoCol.Name { autoIndex = i }
        }

        vals[autoIndex] = tbl.GetGId()

        sqlTpl      = "insert into " + shardTable.Name
        if typ == 2 {
            sqlTpl  += " select %s from %s"
            // TODO: need to parse the select other from table shard db.....

            sql     = fmt.Sprintf(sqlTpl, strings.Join(vals, ", "), other)
        } else {
            sqlTpl  += " values(%s)"
            sql     = fmt.Sprintf(sqlTpl, strings.Join(vals, ", "))
        }
    default:
        return "", []interface{}{}, errors.New("Sorry, the sql has problem, please to check it")
    }

    return sql, params, nil
}
