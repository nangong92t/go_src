package controllers

import (
	"std/data-service/std/app/models"
	"database/sql"
	_ "github.com/Go-SQL-Driver/mysql"
	"github.com/coopernurse/gorp"
    r "github.com/revel/revel"
    "github.com/robfig/revel/modules/db/app"
    "fmt"
)

var (
	Dbm *gorp.DbMap
)

func Init() {
	db.Init()
	//MySQL InnoDB UTF8
    fmt.Printf("%#v\n", db.Db)
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

    /*
	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}
    */
	//创建一个User测试表
	//t := Dbm.AddTable(models.User{}).SetKeys(true, "UserId")
	// setColumnSizes(t, map[string]int{
	//	"Name": 20,
	//})

	//Dbm.TraceOn("[gorp]", r.INFO)
	//Dbm.CreateTables()

	//插入一组测试数据
	demoUser := &models.User{0, "Hobo", "123456", 123456}
	if err := Dbm.Insert(demoUser); err != nil {
		panic(err)
	}

}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
