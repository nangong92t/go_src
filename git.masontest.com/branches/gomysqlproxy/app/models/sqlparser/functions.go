// Copyright 2014, Successfulmatch Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>, 
// Build on dev-0.0.1
// MIT Licensed

// parse the sql submodel.

package sqlparser

import (
    "strconv"
    "regexp"
    "strings"
)

const (
    ADD_ALTER     = iota
    DROP_ALTER
    CHANGE_ALTER
    CREATE_INDEX
    DROP_INDEX
)

// The result eg:
func ParseInsertPart(sql, tableName string) (cols, vals []string, other string, insertType int) {
    insertType  = 0

    sql     = strings.ToLower(sql)

    // to parse the full select insert sql
    re, _       := regexp.Compile("["+tableName + "\\s+|"+tableName+"\\s?]\\((.*)\\)\\s+select\\s+(.*)\\s+from\\s+(.*)")
    subMatch    := re.FindSubmatch([]byte(sql))
    if len(subMatch) == 4 {
        cols    = strings.Split(string(subMatch[1]), ",")
        vals    = strings.Split(string(subMatch[2]), ",")
        other   = string(subMatch[3])
        insertType  = 1
        return
    }

    // to parse the select insert sql
    re, _       = regexp.Compile("\\s+select\\s+(.*)\\s+from\\s+(.*)")
    subMatch    = re.FindSubmatch([]byte(sql))
    if len(subMatch) == 3 {
        vals    = strings.Split(string(subMatch[1]), ",")
        cols    = []string{}
        other   = string(subMatch[2])
        insertType  = 2

        for i, one  := range vals { vals[i] = strings.TrimSpace(one) }

        return
    }

    // to parse the full insert sql
    re, _       = regexp.Compile("["+tableName + "\\s+|"+tableName+"\\s?]\\((.*)\\)\\s+values\\((.*)\\)")
    subMatch    = re.FindSubmatch([]byte(sql))
    if len(subMatch) == 3 {
        cols    = strings.Split(string(subMatch[1]), ",")
        vals    = strings.Split(string(subMatch[2]), ",")

        for i, one  := range cols { cols[i] = strings.TrimSpace(one) }
        for i, one  := range vals { vals[i] = strings.TrimSpace(one) }

        insertType  = 3

        return
    }

    // to parse the simple insert sql
    re, _       = regexp.Compile(tableName + "\\s+values\\((.*)\\)")
    subMatch    = re.FindSubmatch([]byte(sql))
    if len(subMatch) == 2 {
        cols    = []string{}
        vals    = strings.Split(string(subMatch[1]), ",")
        insertType  = 4

        for i, one  := range vals { vals[i] = strings.TrimSpace(one) }

        return
    }

    return
}

func splitTypeLeng(typ string) (string, string) {
    typeName    := ""
    typeLen     := ""

    re, _  := regexp.Compile("(.*)\\((.*)\\)")
    subMatch   := re.FindSubmatch([]byte(typ))
    if len(subMatch) == 3 {
        typeName    = string(subMatch[1])
        typeLen     = string(subMatch[2])
    } else {
        typeName    = typ
    }

    return typeName, typeLen
}

func ParseAlterPart(sql, tableName string) map[string]map[string]string {
    moreAction  := strings.Split(sql, ",")

    cols        := map[string]map[string]string{}

    for i, one := range moreAction {
        if i > 0 { one = "alter table " + tableName + " " + one }
        action, col := ParseAlterPartSingle(one, tableName)
        for k, _    := range col {
            col[k]["action_id"] = strconv.Itoa(action)
            cols[k] = col[k]
        }
    }

    return cols
}

func ParseAlterPartSingle(sql, tableName string) (action int, cols map[string]map[string]string) {

    cols    = map[string]map[string]string{}

    // change all sql to lower letter.
    sql     = strings.ToLower(sql)


    // to parse the CREATE_INDEX
    re, _       := regexp.Compile("\\s*create\\s+index\\s+")
    index       := re.FindIndex([]byte(sql))

    if len(index) == 2 {
        action  = CREATE_INDEX
        return
    }
    // to parse the CREATE_INDEX case 2
    re, _       = regexp.Compile("\\s*add\\s+index|unique|primary\\s+key\\s+")
    index       = re.FindIndex([]byte(sql))

    if len(index) == 2 {
        action  = CREATE_INDEX
        return
    }

    // to parse the ADD_ALTER
    re, _       = regexp.Compile("^\\s*alter.*\\s+"+tableName+"\\s+add\\s+(.*)")
    subMatch    := re.FindSubmatch([]byte(sql))
    colName     := ""
    colOption   := map[string]string{}

    if len(subMatch) == 2 {
        action  = ADD_ALTER

        subStrs := strings.Split(string(subMatch[1]), " ")

        foundIdx    := 1

        for _, one  := range subStrs {
            one = strings.TrimSpace(one)
            if one == "" { continue }

            switch foundIdx {
            case 1:
                colName = one
            case 2:
                colOption["type"], colOption["leng"]    = splitTypeLeng(one)
            case 3:
                re2, _  := regexp.Compile("\\((.*)\\)")
                subMatch2   := re2.FindSubmatch([]byte(one))
                if len(subMatch2) == 2 {
                    colOption["leng"]   = string(subMatch2[1])
                } else {
                    colOption[one]  = ""
                }
            default:
                colOption[one]  = ""
            }

            foundIdx++
        }

        cols[colName]   = colOption

        return
    }

    // to parse the DROP_ALTER
    re, _       = regexp.Compile("^\\s*alter.*\\s+"+tableName+"\\s+drop\\s+column\\s+(.*)\\s*")
    subMatch    = re.FindSubmatch([]byte(sql))

    if len(subMatch) == 2 {
        action  = DROP_ALTER
        cols["name"]   = map[string]string{
            "col": string(subMatch[1]),
        }

        return
    }

    // to parse the CHANGE_ALTER for change
    re, _       = regexp.Compile("^\\s*alter.*\\s+"+tableName+"\\s+change\\s+(.*)")
    subMatch    = re.FindSubmatch([]byte(sql))

    if len(subMatch) == 2 {
        action  = CHANGE_ALTER
        subStrs := strings.Split(string(subMatch[1]), " ")

        foundIdx    := 1

        for _, one  := range subStrs {
            one = strings.TrimSpace(one)
            if one == "" { continue }

            switch foundIdx {
            case 1:
                colName = one
            case 2:
                colOption["newName"]    = one
            case 3:
                colOption["type"], colOption["leng"]    = splitTypeLeng(one)
            case 4:
                re2, _  := regexp.Compile("\\((.*)\\)")
                subMatch2   := re2.FindSubmatch([]byte(one))
                if len(subMatch2) == 2 {
                    colOption["leng"]   = string(subMatch2[1])
                } else {
                    colOption[one]  = ""
                }
            default:
                colOption[one]  = ""
            }

            foundIdx++
        }

        cols[colName]   = colOption

        return
    }

    // to parse the CHANGE_ALTER for modify
    re, _       = regexp.Compile("^\\s*alter.*\\s+"+tableName+"\\s+modify\\s+(.*)")
    subMatch    = re.FindSubmatch([]byte(sql))

    if len(subMatch) == 2 {
        action  = CHANGE_ALTER
        subStrs := strings.Split(string(subMatch[1]), " ")

        foundIdx    := 1

        for _, one  := range subStrs {
            one = strings.TrimSpace(one)
            if one == "" || one == "column" { continue }

            switch foundIdx {
            case 1:
                colName = one
            case 2:
                colOption["type"], colOption["leng"]    = splitTypeLeng(one)
            case 3:
                re2, _  := regexp.Compile("\\((.*)\\)")
                subMatch2   := re2.FindSubmatch([]byte(one))
                if len(subMatch2) == 2 {
                    colOption["leng"]   = string(subMatch2[1])
                } else {
                    colOption[one]  = ""
                }
            default:
                colOption[one]  = ""
            }

            foundIdx++
        }

        cols[colName]   = colOption

        return

    }

    return
}
