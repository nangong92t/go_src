package tabletserver

import (
	"time"

	"github.com/youtube/vitess/go/hack"
	mproto "github.com/youtube/vitess/go/mysql/proto"
	"github.com/youtube/vitess/go/sqltypes"
	"github.com/youtube/vitess/go/vt/context"
	"github.com/youtube/vitess/go/vt/dbconnpool"
	"github.com/youtube/vitess/go/vt/sqlparser"
)

type RequestContext struct {
	ctx      context.Context
	logStats *SQLQueryStats
	qe       *QueryEngine
}

func (rqc *RequestContext) qFetch(logStats *SQLQueryStats, parsedQuery *sqlparser.ParsedQuery, bindVars map[string]interface{}, listVars []sqltypes.Value) (result *mproto.QueryResult) {
	sql := rqc.generateFinalSql(parsedQuery, bindVars, listVars, nil)
	q, ok := rqc.qe.consolidator.Create(string(sql))
	if ok {
		defer q.Broadcast()
		waitingForConnectionStart := time.Now()
		conn, err := rqc.qe.connPool.Get()
		logStats.WaitingForConnection += time.Now().Sub(waitingForConnectionStart)
		if err != nil {
			q.Err = NewTabletErrorSql(FATAL, err)
		} else {
			defer conn.Recycle()
			q.Result, q.Err = rqc.execSQLNoPanic(conn, sql, false)
		}
	} else {
		logStats.QuerySources |= QUERY_SOURCE_CONSOLIDATOR
		q.Wait()
	}
	if q.Err != nil {
		panic(q.Err)
	}
	return q.Result
}

func (rqc *RequestContext) directFetch(conn dbconnpool.PoolConnection, parsedQuery *sqlparser.ParsedQuery, bindVars map[string]interface{}, listVars []sqltypes.Value, buildStreamComment []byte) (result *mproto.QueryResult) {
	sql := rqc.generateFinalSql(parsedQuery, bindVars, listVars, buildStreamComment)
	return rqc.execSQL(conn, sql, false)
}

// fullFetch also fetches field info
func (rqc *RequestContext) fullFetch(conn dbconnpool.PoolConnection, parsedQuery *sqlparser.ParsedQuery, bindVars map[string]interface{}, listVars []sqltypes.Value, buildStreamComment []byte) (result *mproto.QueryResult) {
	sql := rqc.generateFinalSql(parsedQuery, bindVars, listVars, buildStreamComment)
	return rqc.execSQL(conn, sql, true)
}

func (rqc *RequestContext) fullStreamFetch(conn dbconnpool.PoolConnection, parsedQuery *sqlparser.ParsedQuery, bindVars map[string]interface{}, listVars []sqltypes.Value, buildStreamComment []byte, callback func(*mproto.QueryResult) error) {
	sql := rqc.generateFinalSql(parsedQuery, bindVars, listVars, buildStreamComment)
	rqc.execStreamSQL(conn, sql, callback)
}

func (rqc *RequestContext) generateFinalSql(parsedQuery *sqlparser.ParsedQuery, bindVars map[string]interface{}, listVars []sqltypes.Value, buildStreamComment []byte) string {
	bindVars[MAX_RESULT_NAME] = rqc.qe.maxResultSize.Get() + 1
	sql, err := parsedQuery.GenerateQuery(bindVars, listVars)
	if err != nil {
		panic(NewTabletError(FAIL, "%s", err))
	}
	if buildStreamComment != nil {
		sql = append(sql, buildStreamComment...)
	}
	// undo hack done by stripTrailing
	sql = restoreTrailing(sql, bindVars)
	return hack.String(sql)
}

func (rqc *RequestContext) execSQL(conn dbconnpool.PoolConnection, sql string, wantfields bool) *mproto.QueryResult {
	result, err := rqc.execSQLNoPanic(conn, sql, true)
	if err != nil {
		panic(err)
	}
	return result
}

func (rqc *RequestContext) execSQLNoPanic(conn dbconnpool.PoolConnection, sql string, wantfields bool) (*mproto.QueryResult, error) {
	if timeout := rqc.qe.queryTimeout.Get(); timeout != 0 {
		qd := rqc.qe.connKiller.SetDeadline(conn.Id(), time.Now().Add(timeout))
		defer qd.Done()
	}

	start := time.Now()
	result, err := conn.ExecuteFetch(sql, int(rqc.qe.maxResultSize.Get()), wantfields)
	rqc.logStats.AddRewrittenSql(sql, start)
	if err != nil {
		return nil, NewTabletErrorSql(FAIL, err)
	}
	return result, nil
}

func (rqc *RequestContext) execStreamSQL(conn dbconnpool.PoolConnection, sql string, callback func(*mproto.QueryResult) error) {
	start := time.Now()
	err := conn.ExecuteStreamFetch(sql, callback, int(rqc.qe.streamBufferSize.Get()))
	rqc.logStats.AddRewrittenSql(sql, start)
	if err != nil {
		panic(NewTabletErrorSql(FAIL, err))
	}
}
