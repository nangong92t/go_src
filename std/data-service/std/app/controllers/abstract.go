package controllers

import (
    "errors"
    "github.com/coopernurse/gorp"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "std/data-service/std/app/models"
    "std/data-service/std/app/modules"
    "git.masontest.com/branches/goserver/workers/revel"
    "github.com/garyburd/redigo/redis"
)

var (
    // 数据库连接对象
    Dbm *gorp.DbMap
    RedisConn redis.Conn
)

// 抽象控制器, 所有具体实现控制器必须继承于此.
type AbstractController struct {
    *revel.Controller
    //Db *gorp.Transaction
    Db *gorp.DbMap
    Sess *modules.Session
    Rc redis.Conn
}

// 初始化数据库连接
var InitDb func() = func() {
    connectionString := getConnectionString()
    if db, err := sql.Open("mysql", connectionString); err != nil {
        panic(err)
    } else {
        Dbm = &gorp.DbMap{
            Db: db,
            Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
        }
    }
    models.InitDbTableMap(Dbm)

    // 初始化Redis连接.
    prot, hostStr   := getRedisConnString()
    conn, err := redis.Dial(prot, hostStr)
    if err != nil {
        panic(err)
    } else {
        RedisConn   = conn
    }
}

func (c *AbstractController) Begin() revel.Result {
    InitDb()
    // txn, err := Dbm.Begin()

    c.Db    = Dbm
    c.Rc    = RedisConn
    c.Sess  = modules.NewSession(c.Db)

    return nil
}

func (c *AbstractController) Commit() revel.Result {
    if c.Db != nil {
        c.Db.Db.Close()
        c.Db = nil
    }

    c.Rc.Close()
    c.Rc = nil

    return nil
}

func (c *AbstractController) Rollback() revel.Result {
    if c.Db == nil {
        return nil
    }

    c.Db = nil
    return nil
}

func (c *AbstractController) GetUserId() int64 {
    if c.Sess.Session.Id == 0 { return 0 }
    return c.Sess.Session.SessionVal.UserId
}

// 统一渲染未登陆输出.
func (c *AbstractController) RenderLogin() revel.Result {
    return c.Render(nil, errors.New("Please login first."))
}

// 统一渲染无管理权限错误.
func (c *AbstractController) RenderNoAdmin() revel.Result {
    return c.Render(nil, errors.New("No permissions."))
}

