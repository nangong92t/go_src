package plugins

import (
    "strings"
    "os"
)

func ParseArg(args []string) map[string]string {
    result  :=  make(map[string]string)

    for i:=0; i<len(args); i++ {
        oneKV   := strings.Split(args[i], "=")
        if len(oneKV) > 1 {
            result[ strings.Trim(oneKV[0], " ") ] = strings.Trim(oneKV[1], " ")
        }
    }

    return result
}

func IsInArg(name string) bool {
    args    := ParseArg(os.Args)
    return args[name] != ""
}

func GetArg(name string) string {
    args    := ParseArg(os.Args)
    return args[name]
}

func MakeFolder(rootPath string, fileName string) error {
    tPath           := rootPath + "/" + fileName
    isExists, err   := IsFolder(tPath)
    if err != nil { return err }
    if !isExists {
        os.MkdirAll(tPath, 0777)
    }

    return nil
}

// 检测文件夹是否存在.
func IsFolder(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

// check the file is exists.
func IsExistsFile(filename string) bool {
    exists  := true
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        exists  = false
    }
    return exists
}


