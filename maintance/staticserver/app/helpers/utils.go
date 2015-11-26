// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package helpers

import (
    "fmt"
    "crypto/md5"
    "encoding/hex"
    "encoding/gob"
    "bytes"
    "labix.org/v2/mgo/bson"
    "strconv"
    "strings"
    "regexp"
    "math"
)

// 查找字符串数组中是否有指定字符串
func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// 查找int64数组中是否有指定值.
func Int64InSlice(a int64, list []int64) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// 查找ObjectId数组中是否有指定值.
func ObjectIdInSlice(id bson.ObjectId, list []bson.ObjectId) bool {
    for _, b := range list {
        if b.Hex() == id.Hex() {
            return true
        }
    }
    return false
}

// 过滤出Object Id数组中重复的数据.
func FilterOutRepeatObjectId(list []bson.ObjectId) []bson.ObjectId {
    listMap := map[bson.ObjectId]int{}
    result  := []bson.ObjectId{}
    for i:=0; i<len(list); i++ {
        if _, isOk := listMap[list[i]]; !isOk {
            listMap[list[i]]    = 1;
            result  = append(result, list[i])
        }
    }

    return result
}

// to restore into string array for the one column more data.
func RestoreMoreDataInOneColumn(value int64, options map[int]string) []string  {
    f   := strconv.FormatInt(value, 2)

    fs  := strings.Split(f, "")
    fsLen   := len(fs)
    res := []string{}
    j   := -1;
    for i:=fsLen-1; i>=0; i-- {
        j++
        if fs[i] == "0" { continue }
        curI    := int( math.Pow(2, float64(j)) )
        if options[curI] != "" {
            res = append(res, options[curI])
        }
    }

    return res
}

// get ids array from "1,2,3,4"
func GetArrayFromString(s string) []interface{} {
    id     := []interface{}{}
    ids    := strings.Split(s, ",")

    for _, one := range ids {
        if one == "" { continue }
        id = append(id, one)
    }

    return id
}

// get ids int64 array to interface array
func GetInterfaceArrayFromInt64(ids []int64) []interface{} {
    id     := []interface{}{}

    for _, one := range ids {
        if one == 0 { continue }
        id = append(id, one)
    }

    return id
}

func Md5(str string) string {
    h := md5.New()
    h.Write([]byte(str))

    return fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
}

func VerifyEmail(email string) bool {
    var emailPattern = regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")

    return emailPattern.Match([]byte(email))
}

//
// 截取字符串一定长度子字符串.
//
func Substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0

    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length

    if start > end {
        start, end = end, start
    }

    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }

    return string(rs[start:end])
}

func GetBytes(key interface{}) ([]byte, error) {
    var buf bytes.Buffer
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}
