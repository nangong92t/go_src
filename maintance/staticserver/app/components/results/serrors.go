// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// 具体错误定义

package results

import (
    "fmt"
    "strconv"
    "strings"
)

var (
    RenderErrorsMap  = map[int]*RenderErrorResult{}
    DontKnowErr     = &RenderErrorResult{Error: "The error of ignorance", ErrorCode:1}

    RenderErrors = []RenderErrorResult{
        // ------ 系统级错误 ------
        {ErrorCode: 10001, Error: "System error"},
        {ErrorCode: 10002, Error: "Service unavailable"},
        {ErrorCode: 10003, Error: "Remote service error"},
        {ErrorCode: 10004, Error: "IP limit"},
        {ErrorCode: 10006, Error: "Source paramter (appkey) is missing"},
        {ErrorCode: 10007, Error: "Miss required parameter (%s) , see doc for more info"},
        {ErrorCode: 10008, Error: "Too many pending tasks, system is busy"},
        {ErrorCode: 10009, Error: "Job expired"},
        {ErrorCode: 10010, Error: "Illegal request"},
        {ErrorCode: 10012, Error: "Request api not found"},
        {ErrorCode: 10015, Error: "No upload root path"},
        {ErrorCode: 10016, Error: "Directory listing not allowed"},

        // ------ 业务级错误 ------
        //  \---- Auth ----
        {ErrorCode: 20101, Error: "Invalid grant"},
        {ErrorCode: 20102, Error: "Auth faild"},
        {ErrorCode: 20103, Error: "Username or password error"},
        {ErrorCode: 20104, Error: "Invalid app client"},
        {ErrorCode: 20105, Error: "User does not exists"},
        {ErrorCode: 20106, Error: "Unsupported image type, only suport JPG, GIF, PNG"},
        {ErrorCode: 20107, Error: "Image size too large"},
        {ErrorCode: 20108, Error: "Content is null"},
        {ErrorCode: 20109, Error: "Does multipart has image"},
        {ErrorCode: 20110, Error: "Can not write the file to system"},
        {ErrorCode: 20111, Error: "Token invalid"},
        {ErrorCode: 20112, Error: "Permission denied, need a high level user role"},
        {ErrorCode: 20113, Error: "Invalid parameter value, see doc for more info"},
        {ErrorCode: 20114, Error: "Refresh token invalid"},
        {ErrorCode: 20115, Error: "Content too long, please input text less than 140 characters"},

        //  \----Permission ----
        {ErrorCode: 20201, Error: "Permission opration name was empty, please input text"},
        {ErrorCode: 20202, Error: "Permission opration name too long, please input text less than 40 characters"},
        {ErrorCode: 20203, Error: "Role name was empty, please input text"},
        {ErrorCode: 20204, Error: "Role does not exist"},
        {ErrorCode: 20205, Error: "Can not change role name: \"admin\", \"member\", \"gold member\""},
        {ErrorCode: 20206, Error: "Permission does not exists"},

        //  \----User and Account ----
        {ErrorCode: 20301, Error: "The account was exists"},
        {ErrorCode: 20302, Error: "Invalid email"},
        {ErrorCode: 20303, Error: "Email is null"},
        {ErrorCode: 20304, Error: "Password is null"},
        {ErrorCode: 20305, Error: "Password less than six characters long"},
        {ErrorCode: 20306, Error: "Invalid member"},
        {ErrorCode: 20307, Error: "Current user no any role"},
        {ErrorCode: 20308, Error: "Invalid gold member"},
        {ErrorCode: 20309, Error: "Invalid profile"},
        {ErrorCode: 20310, Error: "Display name is null"},
        {ErrorCode: 20311, Error: "Display less than two characters long"},
        {ErrorCode: 20312, Error: "Gender is null"},
        {ErrorCode: 20313, Error: "Avator is null"},
        {ErrorCode: 20314, Error: "Id invalid"},
        {ErrorCode: 20315, Error: "Invalid user"},
        {ErrorCode: 20316, Error: "The name has been repeated"},
        {ErrorCode: 20317, Error: "The data not exists"},
        {ErrorCode: 20318, Error: "Invalid device"},

        // ------- Interest part -------
        {ErrorCode: 20401, Error: "Interest does not exists"},
        {ErrorCode: 20402, Error: "Invalid like type"},
        {ErrorCode: 20403, Error: "Has liked"},
        {ErrorCode: 20404, Error: "Invalid comment type"},

        // ------- Plus or Topic part -------
        {ErrorCode: 20501, Error: "Plus does not exists"},

        // ------- User Backup data part -------
        {ErrorCode: 20601, Error: "Invalid backup data type"},
    }
)

func init() {
    return
    for i:=0; i<len(RenderErrors); i++ {
        one := RenderErrors[i]
        RenderErrorsMap[ one.ErrorCode ] = &one
    }
}

// 统一返回错误定义
func GetError(ids []interface{}) []*RenderErrorResult {
    curErr  := []*RenderErrorResult{}
    for i:=0; i<len(ids); i++ {
        curId   := ids[i]
        switch curId.(type) {
        case int:
            if curErrStr, isOk := RenderErrorsMap[curId.(int)]; isOk {
                curErr  = append(curErr, curErrStr)
            } else {
                DontKnowErr.ErrorCode   = curId.(int)
                curErr  = append(curErr, DontKnowErr)
            }
        case error:
            err     := curId.(error)
            errStr  := fmt.Sprintf("%s", err)
            errId, e:= strconv.Atoi(errStr)
            if e == nil {
                if curE, isOk := RenderErrorsMap[errId]; isOk {
                    curErr  = append(curErr, curE)
                } else {
                    DontKnowErr.ErrorCode   = 1
                    curErr  = append(curErr, DontKnowErr)
                }

            } else if errIds := strings.Split(errStr, "[&&]"); len(errIds) > 0 {
                for j:=0; j<len(errIds); j++ {
                    id, e := strconv.Atoi(errIds[j])
                    if e!= nil {
                        fmt.Printf("system error 1: %#v\n", err)
                        curErr  = append(curErr, RenderErrorsMap[10001])
                        continue
                    }
                    if curE, isOk := RenderErrorsMap[id]; isOk {
                        curErr  = append(curErr, curE)
                    }
                }
            } else {
                fmt.Printf("system error 2: %#v\n", err)
                curErr  = append(curErr, RenderErrorsMap[10001])
            }
        }
    }

    return curErr
}

