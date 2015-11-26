// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// to prepare the select execute planning.

package planbuilder

import (
	"fmt"
    "log"
    "errors"
    "runtime"
    "strings"
    "strconv"
    "sync"

    "git.masontest.com/branches/gomysqlproxy/app/models/host"
    "git.masontest.com/branches/gomysqlproxy/app/models/schema"
    "git.masontest.com/branches/gomysqlproxy/app/models/client"
    "git.masontest.com/branches/gomysqlproxy/app/models/sqlparser"
)

type SelectDML struct {
    DML
    PDML        *sqlparser.Select
}

type selectResult struct {
    Res     []map[string]interface{}
    Err     error
}

func NewSelectDML(sel *sqlparser.Select, tableName string, sql []string, user *client.Client) *SelectDML {
    d   := &SelectDML {}
    d.SetAction("select")

    d.PDML  = sel
    d.TableName = tableName
    d.Sql   = sql
    d.User  = user
    d.Mut       = &sync.Mutex{}

    return d
}

func analyzeSelect(sel *sqlparser.Select, args []string, user *client.Client) (plan *ExecPlan, err error) {
	// Default plan
	plan = &ExecPlan{
		PlanId:     PLAN_PASS_SELECT,
		FieldQuery: GenerateFieldQuery(sel),
		FullQuery:  GenerateSelectLimitQuery(sel),
	}

	// from
	tableName, _ := analyzeFrom(sel.From)
	if tableName == "" {
		plan.Reason = REASON_TABLE
		return plan, nil
	}

    // init the plan core object.
    plan.Core       = NewSelectDML(sel, tableName, args, user)

	// There are bind variables in the SELECT list
	if plan.FieldQuery == nil {
		plan.Reason = REASON_SELECT_LIST
		return plan, nil
	}

	if sel.Distinct != "" || sel.GroupBy != nil || sel.Having != nil {
		plan.Reason = REASON_SELECT
		return plan, nil
	}

	return plan, nil
}


func analyzeSelectExprs(exprs sqlparser.SelectExprs, table *schema.MysqlTable) (selects []int, err error) {
	selects = make([]int, 0, len(exprs))
	for _, expr := range exprs {
		switch expr := expr.(type) {
		case *sqlparser.StarExpr:
			// Append all columns.
			for colIndex := range table.Columns {
				selects = append(selects, colIndex)
			}
		case *sqlparser.NonStarExpr:
			name := sqlparser.GetColName(expr.Expr)
			if name == "" {
				// Not a simple column name.
				return nil, nil
			}
			colIndex := table.FindColumn(name)
			if colIndex == -1 {
				return nil, fmt.Errorf("column %s not found in table %s", name, table.Name)
			}
			selects = append(selects, colIndex)
		default:
			log.Println("unreachable")
		}
	}
	return selects, nil
}


func analyzeFrom(tableExprs sqlparser.TableExprs) (tablename string, hasHints bool) {
	if len(tableExprs) > 1 {
		return "", false
	}
	node, ok := tableExprs[0].(*sqlparser.AliasedTableExpr)
	if !ok {
		return "", false
	}
	return sqlparser.GetTableName(node.Expr), node.Hints != nil
}

func analyzeWhere(node *sqlparser.Where) (conditions []sqlparser.BoolExpr) {
	if node == nil {
		return nil
	}
	return analyzeBoolean(node.Expr)
}

func analyzeBoolean(node sqlparser.BoolExpr) (conditions []sqlparser.BoolExpr) {
	switch node := node.(type) {
	case *sqlparser.AndExpr:
		left := analyzeBoolean(node.Left)
		right := analyzeBoolean(node.Right)
		if left == nil || right == nil {
			return nil
		}
		if sqlparser.HasINClause(left) && sqlparser.HasINClause(right) {
			return nil
		}
		return append(left, right...)
	case *sqlparser.ParenBoolExpr:
		return analyzeBoolean(node.Expr)
	case *sqlparser.ComparisonExpr:
		switch {
		case sqlparser.StringIn(
			node.Operator,
			sqlparser.AST_EQ,
			sqlparser.AST_LT,
			sqlparser.AST_GT,
			sqlparser.AST_LE,
			sqlparser.AST_GE,
			sqlparser.AST_NSE,
			sqlparser.AST_LIKE):
			if sqlparser.IsColName(node.Left) && sqlparser.IsValue(node.Right) {
				return []sqlparser.BoolExpr{node}
			}
		case node.Operator == sqlparser.AST_IN:
			if sqlparser.IsColName(node.Left) && sqlparser.IsSimpleTuple(node.Right) {
				return []sqlparser.BoolExpr{node}
			}
		}
	case *sqlparser.RangeCond:
		if node.Operator != sqlparser.AST_BETWEEN {
			return nil
		}
		if sqlparser.IsColName(node.Left) && sqlparser.IsValue(node.From) && sqlparser.IsValue(node.To) {
			return []sqlparser.BoolExpr{node}
		}
	}
	return nil
}


// to come true ExecCore interface
func (d *SelectDML) Destory() {
    d.Hosts = nil
    d.User  = nil
    d.PDML  = nil
}

// to check is wheather have desc in order by.
func (d *SelectDML) hasDescInOrder() bool {
    result  := false
    if d.PDML.OrderBy != nil {
        for _, one := range d.PDML.OrderBy {
            if one.Direction == "desc" {
                result  = true
                break
            }
        }
    }
    return result
}

// to come true ExecCore interface
func (d *SelectDML) Do() (interface{}, error) {
    d.Mut.Lock()
    defer d.Mut.Unlock()

    var err error
    var shardTables []*schema.MysqlShardTable
    var curTable    *schema.MysqlTable

    for i:=0; i<len(schema.Tables); i++ {
        if schema.Tables[i].Name == d.TableName {
            curTable    = schema.Tables[i]
            shardTables = curTable.Shards
            break
        }
    }
    if shardTables == nil { return nil, errors.New("no found any one this table " + d.TableName) }

	// Select expressions
	selects, err    := analyzeSelectExprs(d.PDML.SelectExprs, curTable)
    limit           := 0

    if d.PDML.Limit != nil {
        // offset, err  = strconv.Atoi(sqlparser.String(d.PDML.Limit.Offset))
        // if err != nil { return nil, err }

        limit, err   = strconv.Atoi(sqlparser.String(d.PDML.Limit.Rowcount))
        if err != nil { return nil, err }
    }

	if err != nil {
		return nil, err
	}
	if selects == nil {
        return nil, errors.New("no any found select column")
	}

	// where
    // conditions  := analyzeWhere(d.PDML.Where)
    params      := getInterfaceParams(d.Sql[1:])

    resultCh    := make(chan *selectResult, runtime.NumCPU())
    shardTBLen  := len(shardTables)
    maxTotal    := 5000


    // to select the data from all shard tabl at the same time.
    for i:=0; i < shardTBLen; i++ {
        go func(i int) {
            stb := shardTables[i]
            conn, err := host.GetBetterHost(stb.ShardDB.HostGroup.Slave, "slave").ConnToDB(stb.ShardDB.Name)
            if err != nil {
                conn.Close()
                resultCh <- &selectResult{Res:nil, Err: err}
            }

            trueSql     := strings.Replace(d.Sql[0], d.TableName, stb.Name, -1)
            if limit == 0 { trueSql += " limit " + strconv.Itoa(maxTotal) }

            stmt, err   := conn.Prepare(trueSql)
            if err != nil {
                conn.Close()
                resultCh <- &selectResult{Res:nil, Err: err}
            }
            rows, err    := stmt.Query(params...)
            if err != nil {
                conn.Close()
                resultCh <- &selectResult{Res:nil, Err: err}
            }

            stmt.Close()
            conn.Close()

            dataCols, err   := rows.Columns()
            if err != nil {
                resultCh <- &selectResult{Res:nil, Err: err}
            }

            dataLen := len(dataCols)
            rowData := make([]interface{}, dataLen)
            valuePtrs := make([]interface{}, dataLen)
            backData    := []map[string]interface{}{}
            curFound    := 0

            for rows.Next() {
                if curFound >= maxTotal { break }

                oneData := map[string]interface{}{}
                for j, _ := range dataCols {
                    valuePtrs[j] = &rowData[j]
                }

                err = rows.Scan(valuePtrs...)
                if err != nil { resultCh <- &selectResult{Res:nil, Err: err}; return }

                for j, col := range dataCols {
                    var v interface{}
                    val := rowData[j]
                    b, ok := val.([]byte)
                    if (ok) {
                        v = string(b)
                    } else {
                        v = val
                    }
                    oneData[col]    = v
                }

                backData    = append(backData, oneData)
                curFound++
            }
            err = rows.Err()
            rows.Close()

            if err != nil {
                resultCh <- &selectResult{Res:nil, Err: err}
            }
            resultCh <- &selectResult{Res:backData, Err: nil}
        }(i)
    }

    result  := make([][]map[string]interface{}, shardTBLen)

    // get all shard data.
    for i:=0; i < shardTBLen; i++ {
        res := <-resultCh
        if res.Err != nil { return nil, res.Err }

        result[i]  = res.Res
    }

    // to merge all data from all shard.
    mergedResult    := []map[string]interface{}{}

    if d.hasDescInOrder() {           // for desc order data in shard tables
        for i:=(shardTBLen-1); i >= 0; i-- {
            if i==(shardTBLen-1) { mergedResult = result[i]; continue }

            for j:=0; j<len(result[i]); j++ {
                mergedResult    = append(mergedResult, result[i][j])
                if limit != 0 && limit <= len(mergedResult) {
                    goto LAST
                }
            }
        }
    } else {                            // for asc order data in shard tables
        for i:=0; i < shardTBLen; i++ {
            if i==0 { mergedResult = result[i]; continue }

            for j:=0; j<len(result[i]); j++ {
                mergedResult    = append(mergedResult, result[i][j])
                if limit != 0 && limit <= len(mergedResult) {
                    goto LAST
                }
            }
        }
    }

LAST:
    return mergedResult, nil
}
