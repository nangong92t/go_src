package helpers

import (
    "strconv"
    "strings"
    "math"
    "std/data-service/std/app/models"
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

// to merge the date data to the array.
func MergeDateData(source []models.DateData, days, toDate int64) []interface{} {
    result  := []interface{}{}

    dataMap     := map[int64]int64{}
    for _, one := range source {
        dataMap[one.Date]   = one.Val.(int64)
    }

    for i:=int64(0); i<days; i++ {
        curDay  := (toDate - i * 86400) * 1000
        if dataMap[curDay] != 0 {
            result  = append(result, [2]interface{}{
                curDay,
                dataMap[curDay],
            })
        } else {
            result  = append(result, [2]interface{}{
                curDay,
                0,
            })
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
