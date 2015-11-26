// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package models

import (
    "fmt"
    "errors"
    "encoding/json"
    "strings"
    "labix.org/v2/mgo"
    "labix.org/v2/mgo/bson"
)

type M bson.M

// cache one database session.
var (
    shardDbSession *mgo.Session
    EmptyObjectID = bson.ObjectIdHex("000000000000000000000000")
)

// global database model fances.
type DBModel interface {
    GetCollectionName() string
    GetDbName() string

    BeforeSave() (error, bool)
    AfterSave()

    BeforeFind()
    AfterFind()

    BeforeUpdate() (error, bool)
    AfterUpdate()

    BeforeDelete()
    AfterDelete()

    SetDefault()
    Validate() error

    GetId() bson.ObjectId
}

type Gps struct {
    Lat             float32         `json:"lat" bson:"lat"`
    Lon             float32         `json:"lon" bson:"lon"`
}

type AbstractModel struct {
    DBModel                     `json:"-"`

    // to connect inheritance object for come back from parent object.
    Self        DBModel         `json:"-"`

    // database session
    session     *mgo.Session    `bson:"-", json:"-"`

    // is new recode.
    IsNew       bool            `json:"-" bson:"-"`
}

func (a *AbstractModel) GetID() bson.ObjectId {
    return bson.NewObjectId()
}

func (a *AbstractModel) Init(self DBModel) *AbstractModel {
    return nil
    a.getDBSession()
    a.Self  = self

    return a
}

func (a *AbstractModel) getDBSession() {
    session, err := mgo.Dial("127.0.0.1:27017")
    if err != nil { panic(err) }
    a.session   = session
    if shardDbSession == nil {
        shardDbSession = session
    }
}

func (a *AbstractModel) GetDb() *mgo.Session {
    curSession  := a.session
    if curSession == nil {
        if shardDbSession != nil {
            // get cache db session connection.
            curSession = shardDbSession
        } else {
            // init new object session
            a.getDBSession()
            curSession  = a.session
        }
    }

    return curSession.Clone()
}

// abstract function
func (a *AbstractModel) GetDbName() string {
    // default the database name.
    return "maintance/staticserver"
}

// abstract function
func (a *AbstractModel) GetCollectionName() string {
    return ""
}

// abstract function
func (a *AbstractModel) BeforeSave() (error, bool) {
    return a.Self.BeforeSave()
}

// abstract function
func (a *AbstractModel) AfterSave() {
    a.Self.AfterSave()
}

// abstract function
func (a *AbstractModel) GetId() bson.ObjectId { return a.Self.GetId() }
func (a *AbstractModel) BeforeFind() { a.Self.BeforeFind() }
func (a *AbstractModel) AfterFind() { a.Self.AfterFind() }
func (a *AbstractModel) BeforeUpdate() (error, bool) { return a.Self.BeforeUpdate() }
func (a *AbstractModel) AfterUpdate() { a.Self.AfterUpdate() }
func (a *AbstractModel) BeforeDelete() { a.Self.BeforeDelete() }
func (a *AbstractModel) AfterDelete() { a.Self.AfterDelete() }

// abstract function
func (a *AbstractModel) Validate() error {
    return nil
}

func (a *AbstractModel) SetDefault() {}

func (a *AbstractModel) CFind(f func(*mgo.Collection)) {
    dbName  := a.Self.GetDbName()
    cName   := a.Self.GetCollectionName()
    handle  := a.GetDb()
    defer func() {
        handle.Close()
        if err := recover(); err != nil {
            panic(err)
        }
    }()

    c       := handle.DB(dbName).C(cName)

    f(c)
}

func (a *AbstractModel) FindAll(condition M, sort string, skip, limit int, needTotal bool, result interface{}) (total int) {
    var err error

    a.CFind(func (c *mgo.Collection) {
        finder  := c.Find(condition)

        if needTotal { total, _    = finder.Count() }

        if sort != "" { finder.Sort(sort) }
        if skip > 0 { finder.Skip(skip) }
        if limit > 0 { finder.Limit(limit) } else { finder.Limit(5000) }

        var explainResult interface{}
        finder.Explain(&explainResult)
        fmt.Printf("cond: %#v, Explain: %#v\n--------\n", condition, explainResult)

        err     = finder.All(result)
        if err != nil { fmt.Printf("finder have error: %#v\n", err) }

    })

    return
}

func (a *AbstractModel) Find(condition M, result interface{}) (error) {
    var err error

    a.Self.BeforeFind()

    a.CFind(func (c *mgo.Collection) {
        finder  := c.Find(condition)
        err     =  finder.One(result)

        // 当上面find后 a.Self将会丢失， 很奇怪的地方.
        a.Self  = result.(DBModel)
    })

    if err == nil { a.Self.AfterFind() }

    return err
}

func (a *AbstractModel) FindByPk(id bson.ObjectId, result interface{}) (error) {
    return a.Find(M{"_id": id}, result)
}

func (a *AbstractModel) FindByPkStr(id string, result interface{}) (error) {
    if !bson.IsObjectIdHex(id) {
        return errors.New("20314")
    }
    return a.Find(M{"_id": bson.ObjectIdHex(id)}, result)
}

func (a *AbstractModel) Save(validData bool) (err error) {
    var isOk bool

    if validData {
        err = a.Self.Validate()
        if err != nil { return }
    }

    if a.IsNew {
        if err, isOk = a.BeforeSave(); !isOk { return }
    } else {
        if err, isOk = a.BeforeUpdate(); !isOk { return }
    }

    dbName  := a.Self.GetDbName()
    cName   := a.Self.GetCollectionName()
    handle  := a.GetDb()
    defer func() {
        handle.Close()
        if err2 := recover(); err2 != nil {
            switch err2.(type) {
            case string:
                err = errors.New(err2.(string))
            case error:
                err = err2.(error)
            }
        }
    }()

    c       := handle.DB(dbName).C(cName)

    if a.IsNew {
        err = c.Insert(a.Self)
        if err == nil { a.AfterSave() }
    } else {
        err = c.UpdateId(a.Self.GetId(), a.Self)
        if err == nil { a.AfterUpdate() }
    }

    a.IsNew = false

    return
}

func (a *AbstractModel) IsEmptyObjectId(id bson.ObjectId) bool {
    return id.Hex() == "000000000000000000000000"
}

func (a *AbstractModel) Update(selector interface{}, update interface{}, isValid bool) (err error) {
    var isOk bool

    if isValid {
        if err, isOk = a.BeforeUpdate(); !isOk || err != nil { return }
    }

    dbName  := a.Self.GetDbName()
    cName   := a.Self.GetCollectionName()
    handle  := a.GetDb()
    defer func() {
        handle.Close()
        if err2 := recover(); err2 != nil {
            switch err2.(type) {
            case string:
                err = errors.New(err2.(string))
            case error:
                err = err2.(error)
            }
        }
    }()

    c       := handle.DB(dbName).C(cName)

    err     = c.Update(selector, M{"$set": update})

    return
}

func (a *AbstractModel) UpdateByPk(id bson.ObjectId, update interface{}) error {
    if !id.Valid() { return errors.New("20314") }

    return a.Update(M{"_id": id}, update, true)
}

func (a *AbstractModel) Delete(condition bson.M, isAll bool) (err error) {
    dbName  := a.Self.GetDbName()
    cName   := a.Self.GetCollectionName()
    handle  := a.GetDb()
    defer func() {
        handle.Close()

        if err2 := recover(); err2 != nil {
            switch err2.(type) {
            case string:
                err = errors.New(err2.(string))
            case error:
                err = err2.(error)
            }
        }
    }()

    a.Self.BeforeDelete()

    c       := handle.DB(dbName).C(cName)

    if isAll {
        _, err = c.RemoveAll(condition)
    } else {
        err     = c.Remove(condition)
    }

    if err == nil { a.Self.AfterDelete() }

    return err
}

func (a *AbstractModel) DeleteByPk(id bson.ObjectId) error {
    return a.Delete(bson.M{"_id": id}, false)
}

func (a *AbstractModel) DeleteByPkStr(id string) error {
    if !bson.IsObjectIdHex(id) {
        return errors.New("pk id is wrong!")
    }
    return a.Delete(bson.M{"_id": bson.ObjectIdHex(id)}, false)
}

func (a *AbstractModel) JoinStrErrors(errCodes []string) error {
    if len(errCodes) > 0 {
        return errors.New(strings.Join(errCodes, "[&&]"))
    }

    return nil
}

//
// 删除内嵌数组中的id。
//
// @param bson.ObjectId id        主数据id.
// @param bson.ObjectId delId     待删除主Id
// @param string        fieldName 删除的目标字段名.
// @param bool          isAll     是否删除全部，true: 是, false: 否.
//
func (a *AbstractModel) DeleteInArrayById(id, delId bson.ObjectId, fieldName string, isAll bool) (err error) {
    // 其本质是更新.
    dbName  := a.Self.GetDbName()
    cName   := a.Self.GetCollectionName()
    handle  := a.GetDb()
    defer func() {
        handle.Close()
        if err2 := recover(); err2 != nil {
            switch err2.(type) {
            case string:
                err = errors.New(err2.(string))
            case error:
                err = err2.(error)
            }
        }
    }()

    c       := handle.DB(dbName).C(cName)

    if isAll {
        err = c.UpdateId(id, M{"$pullAll": M{fieldName: delId}})
    } else {
        err = c.UpdateId(id, M{"$pull": M{fieldName: delId}})
    }

    return
}

//
// 存在更新， 不存在直接插入.
//
func (a *AbstractModel) Upsert(selector interface{}, update interface{}) (err error) {
    // 其本质是更新.
    dbName  := a.Self.GetDbName()
    cName   := a.Self.GetCollectionName()
    handle  := a.GetDb()
    defer func() {
        handle.Close()
        if err2 := recover(); err2 != nil {
            switch err2.(type) {
            case string:
                err = errors.New(err2.(string))
            case error:
                err = err2.(error)
            }
        }
    }()

    c       := handle.DB(dbName).C(cName)
    _, err  = c.Upsert(selector, update)

    return
}

//
// 转Json字符串.
//
func (a *AbstractModel) ToJson() (string) {
    result, _ := json.Marshal(a.Self);
    return string(result)
}
