// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// 数据层通用工具方法.

package models

import (
    "labix.org/v2/mgo/bson"
)

// 检测ObjectId数组中是否存在某一个id。
func ObjectIdIsIn(ids []bson.ObjectId, id bson.ObjectId) bool {
    result  := false
    for i:=0; i<len(ids); i++ {
        if ids[i].Hex() == id.Hex() {
            result  = true
            break
        }
    }

    return result
}

// 检测并转换object id hex数组为object id数组.
func ToObjectIds(ids []string) []bson.ObjectId {
    res := []bson.ObjectId{}
    for i:=0; i<len(ids); i++ {
        if bson.IsObjectIdHex(ids[i]) {
            res = append(res, bson.ObjectIdHex(ids[i]))
        }
    }

    return res
}
