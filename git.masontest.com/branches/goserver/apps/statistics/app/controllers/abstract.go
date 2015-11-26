package controllers

import (
    "github.com/coopernurse/gorp"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "std/data-service/std/app/models"
    "std/data-service/std/app/modules"
    "github.com/revel/revel"
    "github.com/garyburd/redigo/redis"
    "fmt"
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
        revel.ERROR.Fatal(err)
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

// 统一格式化输出
func (c *AbstractController) Render(data interface{}, err error) revel.Result {
    errorCode   := 0
    errorMsg    := ""

    if err != nil {
        errorCode   = 1
        errorMsg    = fmt.Sprintf("%s", err)
        res := map[string]interface{} {
            "errormsg": errorMsg,
            "errorcode": errorCode,
        }

        return c.RenderJson(res)
    }

    return c.RenderJson(data)
}

// 统一渲染未登陆输出.
func (c *AbstractController) RenderLogin() revel.Result {
    res := map[string]interface{} {
        "errmsg": "Please login first.",
        "errcode": 1,
    }

    return c.RenderJson(res)
}

// 统一渲染无管理权限错误.
func (c *AbstractController) RenderNoAdmin() revel.Result {
    res := map[string]interface{} {
        "errmsg": "No permissions.",
        "errcode": 1,
    }

    return c.RenderJson(res)
}


