// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// profile相关帮助工具.

package helpers

import (
    "time"
    "strconv"
)

// 根据年龄与当前时间获取生日
func GetBirthdayByAge(age int) time.Time {
    year    := time.Now().Year()
    year    = year - age
    birthday, _ := time.Parse("2006-01-02", strconv.Itoa(year) + "-01-01")
    return birthday
}
