// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package modules

import (
    "fmt"
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    "reflect"
    "errors"
    "runtime"
    "labix.org/v2/mgo/bson"

    "maintance/staticserver/app/components/cache/redis"
    "maintance/staticserver/app/helpers"
    //"maintance/staticserver/app/models"
)

type AbstractModule struct {
}

func (a *AbstractModule) Init() {
}

func (a *AbstractModule) GetCacheKey(key string) string {
    h   := md5.New()
    h.Write([]byte(key))

    return fmt.Sprintf("data_cache_%s", hex.EncodeToString(h.Sum(nil)))
}

func (a *AbstractModule) GetFunCacheKey(funcName string, funcParams ... interface{}) string {
    h   := md5.New()

    params  := funcName
    for i:=0; i<len(funcParams); i++ {
        params = params + "_" + fmt.Sprintf("%#v", funcParams[i])
    }
    h.Write([]byte(params))

    return fmt.Sprintf("data_cache_%s", hex.EncodeToString(h.Sum(nil)))
}

func (a *AbstractModule) SetCache(key string, value interface{}, expires int64) error {
    cacheData, _ := json.Marshal(value)
    return redis.Write(key, string(cacheData), expires)
}

func (a *AbstractModule) DelCache(key string) error {
    return redis.Delete(key)
}

func (a *AbstractModule) GetCache(key string) (interface{}, error) {
    data, err         := redis.Read(key)
    if err != nil { return nil, err }
    if data != nil {
        return data, nil
    }

    return nil, err
}

func (a *AbstractModule) VerifyPager(page, limit int) (int, int) {
    if page == 0 { page = 1 }
    if limit == 0 || limit > 5000 { limit = 20 }

    skip    := (page - 1) * limit

    return skip, limit
}

func (a *AbstractModule) BindCacheResult(in interface{}, out interface{}) (err error) {
    inVal   := reflect.ValueOf(in)
    outVal  := reflect.ValueOf(out)

    if outVal.Kind() != reflect.Ptr {
        return errors.New("data argument must be a slice address")
    }
    if inVal.Kind() == reflect.Ptr { inVal  = inVal.Elem() }

    if inVal.Kind() == reflect.Struct {
        outVal  = outVal.Elem()
        outVal.Set( inVal )
        return nil
    } else if inVal.Kind() == reflect.Map {
        inBytes, _ := json.Marshal(in.(map[string]interface{}))
        err = json.Unmarshal(inBytes, out)
    }

    return
}

//
// 通用缓存存取组件方法.
//
// @param getFun    读取数据主方法，当没有数据缓存时才会执行.
// @param result    读取缓存数据的容器struct pointer. 当为nil时，不设置.
// @param expiredIn 多少秒之后过期, 如果为0,代表永久缓存,直到缓存服务不可用为止.
// @param params    当没有数据缓存时， 执行主调用方法时的参数数组，顺序对应主方法接口.
//
func (a *AbstractModule) GetCacheData(getFun interface{}, resData interface{}, expriedIn int64, params ... interface{}) (interface{}, string, error) {
    curFunVal   := reflect.ValueOf(getFun)
    if curFunVal.Kind() != reflect.Func { return nil, "", errors.New("no func") }
    if curFunVal.Type().NumOut() != 2 {
        return nil, "", errors.New("the func must set 2 output params, and first param is main data, last param is error")
    }
    funName     := runtime.FuncForPC(curFunVal.Pointer()).Name()
    key         := a.GetFunCacheKey(funName, params...)

    data, err   := a.GetCache(key)
    if err != nil { return nil, key, err }
    if data != nil {
        if resData != nil {
            err = a.BindCacheResult(data, resData)
            if err != nil { return nil, key, err }
        }
        return data, key, nil
    }

    // 当没有数据时执行主方法，到数据库中获取数据.
    in := make([]reflect.Value, len(params))
    for k, param := range params {
        in[k] = reflect.ValueOf(param)
    }
    result      := curFunVal.Call(in)
    outData     := result[0].Interface()
    outErr      := result[1].Interface()
    if outErr != nil {
        if err, isOk := outErr.(error); isOk {
            return nil, key, err
        } else {
            return nil, key, errors.New("the func last output param must a error type")
        }
    }

    // 当获取到有效数据后， 按照expriedIn时值缓存数据.
    if result[0].IsValid() {
        if result[0].Kind() == reflect.Ptr {
            result[0]   = result[0].Elem()
        }

        if result[0].Kind() == reflect.Slice || result[0].Kind() == reflect.Map || result[0].Kind() == reflect.Array {
            if result[0].Len() > 0 {
                if resData != nil { a.BindCacheResult(data, resData) }
                a.SetCache(key, outData, expriedIn)
            }
        } else if result[0].Kind() == reflect.Struct || result[0].Elem().Kind() == reflect.Struct {
            if resData != nil { a.BindCacheResult(outData, resData) }
            a.SetCache(key, outData, expriedIn)
        }
    }

    return outData, key, nil
}

//
// 通过关键字过滤出数据中唯一主键id数组.
//
func (a *AbstractModule) FilterIdsByField(data interface{}, fieldName string) ([]bson.ObjectId, error) {
    ids := []bson.ObjectId{}
    curVal := reflect.ValueOf(data)
    if curVal.Kind() != reflect.Ptr || curVal.Elem().Kind() != reflect.Slice {
        return ids, errors.New("data argument must be a slice address")
    }

    slicev := curVal.Elem()
    slicev = slicev.Slice(0, slicev.Cap())
    for i:=0; i<slicev.Len(); i++ {
        val := slicev.Index(i).FieldByName(fieldName).Interface()
        if id, isOk := val.(bson.ObjectId); isOk {
            if id.Valid() && !helpers.ObjectIdInSlice(id, ids) {
                ids = append(ids, id)
            }
        } else if idArr, isOk := val.([]bson.ObjectId); isOk {
            for j:=0; j<len(idArr); j++ {
                if idArr[j].Valid() && !helpers.ObjectIdInSlice(idArr[j], ids) {
                    ids = append(ids, idArr[j])
                }
            }
        }
    }

    return ids, nil
}

//
// 合并相关依赖正式数据到对应数据模型中.
//
func (a *AbstractModule) MergeDataTo(data interface{}, fieldName, withDataFieldName string, mapData interface{}) error {
    curVal := reflect.ValueOf(data)
    if curVal.Kind() != reflect.Ptr {
        return errors.New("data argument must be a slice address or struct address")
    }

    mapVal  := reflect.ValueOf(mapData)
    if mapVal.Kind() != reflect.Ptr || mapVal.Elem().Kind() != reflect.Map {
        return errors.New("map argument must be a map address")
    }

    slicev := curVal.Elem()
    mapVal  = mapVal.Elem()

    // 设置一个结构数组
    if slicev.Kind() == reflect.Slice {
        slicev = slicev.Slice(0, slicev.Cap())

        for i:=0; i<slicev.Len(); i++ {
            idVal   := slicev.Index(i).FieldByName(fieldName)
            if !idVal.IsValid() { continue }

            withDataVal := slicev.Index(i).FieldByName(withDataFieldName)
            if !withDataVal.IsValid() { continue }

            // 设置单个数据
            if _, isOk := idVal.Interface().(bson.ObjectId); isOk {
                dataVal := mapVal.MapIndex(idVal)
                if !dataVal.IsValid() || dataVal.IsNil() { continue }

                withDataVal.Set(dataVal)
            } else if idArr, isOk := idVal.Interface().([]bson.ObjectId); isOk {
                // 设置一个数组
                idArrLen:= len(idArr)
                list    := withDataVal.Slice(0, withDataVal.Cap())
                for j:=0; j<idArrLen; j++ {
                    dataVal := mapVal.MapIndex(reflect.ValueOf(idArr[j]))
                    if !dataVal.IsValid() || dataVal.IsNil() { continue }

                    list = reflect.Append(list, dataVal)
                }

                if list.Len() > 0 {
                    withDataVal.Set(list)
                }
            }
        }
    } else {
        // 设置一单个结构
        idVal   := slicev.FieldByName(fieldName)
        if !idVal.IsValid() { return errors.New("no source id field name") }

        withDataVal := slicev.FieldByName(withDataFieldName)
        if !withDataVal.IsValid() { return errors.New("no aim field name") }

        // 设置单个数据
        if _, isOk := idVal.Interface().(bson.ObjectId); isOk {
            dataVal := mapVal.MapIndex(idVal)
            if !dataVal.IsValid() { return nil }

            withDataVal.Set(dataVal)
        } else if idArr, isOk := idVal.Interface().([]bson.ObjectId); isOk {
            // 设置一个数组
            idArrLen:= len(idArr)
            list    := withDataVal.Slice(0, withDataVal.Cap())
            for j:=0; j<idArrLen; j++ {
                dataVal := mapVal.MapIndex(reflect.ValueOf(idArr[j]))
                if !dataVal.IsValid() { continue }

                list = reflect.Append(list, dataVal)
                list = list.Slice(0, list.Cap())
            }

            if list.Len() > 0 {
                withDataVal.Set(list)
            }
        }
    }

    return nil
}
