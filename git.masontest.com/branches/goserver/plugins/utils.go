package plugins

import (
    "strings"
    "os"
    "compress/gzip"
    "io/ioutil"
    "bytes"
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

func GzipEncode(in []byte) ([]byte, error) {
     var (
         buffer bytes.Buffer
         out    []byte
         err    error
     )
     writer := gzip.NewWriter(&buffer)

     _, err = writer.Write(in)
     if err != nil {
         writer.Close()
         return out, err
     }
     err = writer.Close()
     if err != nil {
         return out, err
     }

     return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
     reader, err := gzip.NewReader(bytes.NewReader(in))
     if err != nil {
         var out []byte
         return out, err
     }
     defer reader.Close()

     return ioutil.ReadAll(reader)
}

