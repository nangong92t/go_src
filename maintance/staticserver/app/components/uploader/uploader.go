// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package uploader

import (
    "os"
    "time"
    "fmt"
    "encoding/hex"
    "io/ioutil"
    "crypto/md5"
    "errors"
    "strings"
    "net/http"
    "path/filepath"
    "github.com/revel/revel"
    "maintance/staticserver/app/models"
)

func HelperIsFolderExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

//
// 根据文件名生成文件路径, 如果文件路径不存在并创建
//
// @param bool isRand 是否生成随机目录.
//
func HelperForFilePathCreator(rootPath string, fileName string, isRand bool) (string, string) {
    h               := md5.New()
    if isRand { fileName = fileName + fmt.Sprintf("%s", time.Now().Unix()) }
    h.Write([]byte(fileName))
    fileName        = fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
    pathstrs        := []rune(fileName)
    f1, f2, f3      := string(pathstrs[0]), string(pathstrs[1:4]), string(pathstrs[5:10])

    tPath           := rootPath
    isExists, _     := HelperIsFolderExists(tPath)
    if !isExists {
        os.Mkdir(tPath, 0755)
    }

    tPath           = rootPath + "/" + f1
    isExists, _     = HelperIsFolderExists(tPath)
    if !isExists {
        os.Mkdir(tPath, 0755)
    }

    tPath           = tPath + "/" + f2
    isExists, _     = HelperIsFolderExists(tPath)
    if !isExists {
        os.Mkdir(tPath, 0755)
    }

    tPath           = tPath + "/" + f3
    isExists, _     = HelperIsFolderExists(tPath)
    if !isExists {
        os.Mkdir(tPath, 0755)
    }

    path            := f1 + "/" + f2 + "/" + f3

    return path, fileName
}

func UploadSync(r *http.Request) error {
    //fmt.Printf("Received Sync File\n")
    uploadRoot, err := getUploadRoot()
    if err != nil { return err }

    if !isAllowHost(strings.Split(r.RemoteAddr,":")[0]) { return errors.New("not allow this host") }

    // 获取上传文件数据
    file, handler, err  := r.FormFile("file")
    if err != nil { return err }

    fileName    := handler.Filename
    filePath    := strings.Split(fileName, "/upload/")[1]

    // 检测目录是否存在， 不存在则创建之
    paths   := strings.Split(filePath, "/")
    curPath := uploadRoot + "/" + strings.Join(paths[0:len(paths)-1], "/")

    isExists, _ := HelperIsFolderExists(curPath)
    if !isExists {
        err     = os.MkdirAll(curPath, 0755)
        if err != nil { return err }
    }

    fmt.Printf("Save path: %#v, %#v\n", filePath, uploadRoot)

    data, err           := ioutil.ReadAll(file)
    if err != nil { return err }

    err = ioutil.WriteFile(uploadRoot+"/"+filePath, data, 0644)
    if err != nil { return err }

    //fmt.Printf("Save OK\n")
    return nil
}

func isAllowHost(remoteHost string) bool {
    isAllow := false;
    for i:=0; i<len(StaticHosts); i++ {
        if StaticHosts[i] == remoteHost {
            isAllow = true
            break
        }
    }
    return isAllow
}

func getUploadRoot() (string, error) {
    // 获取配置文件中上传根路径
    revel.Config.SetSection( revel.RunMode )
    uploadRoot, found   := revel.Config.String("uploadpath")
    if !found { return "", errors.New("10015") }

    uploadRoot = uploadRoot + "/upload"

    return uploadRoot, nil
}

func Upload(r *http.Request) (*models.AttachModel, error) {
    uploadRoot, err := getUploadRoot()
    if err != nil { return nil, err }

    // 获取上传文件数据
    file, handler, err  := r.FormFile("file")
    if err != nil { file, handler, err = r.FormFile("Filedata") }
    if err != nil { return nil, nil }

    extension           := strings.TrimLeft(strings.ToLower( filepath.Ext(handler.Filename) ), ".")
    switch extension {
    case "gif", "jpg", "jpeg", "png":
    default:
        return nil, errors.New("20106")
    }

    // 获取文件存储路径
    saveHost            := ""
    savePath, saveName  := HelperForFilePathCreator(uploadRoot, handler.Filename, true)

    data, err           := ioutil.ReadAll(file)
    if err != nil { return nil, err }

    size                := len(data)
    if size > 5242880 {     // 文件容量大小不能超过 5M
        return nil, errors.New("20107")
    }

    // get file hash
    h                   := md5.New()
    h.Write(data)
    hash                := fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))

    // to check is whether same file. 如果存在相同的，就直接返回.
    model               := models.NewAttach()
    //old                 := &models.AttachModel{}
    //model.Find(models.M{"hash":hash}, old)
    //if old.Id.Valid() { return old, nil }

    // 不存在相同的文件，就创建.
    err = ioutil.WriteFile(uploadRoot+"/"+savePath+"/"+saveName+"."+extension, data, 0644)
    if err != nil { return nil, errors.New("20110") }

    // 存储图片数据到库中
    attach  := &models.AttachModel{
        Id:         model.GetID(),
        Name:       handler.Filename,
        Size:       size,
        Ext:        extension,
        Hash:       hash,
        SaveHost:   saveHost,
        SavePath:   savePath,
        SaveName:   saveName+"."+extension,
    }

    // 当有图片集群存在时，需要自动异步同步当前上传图片到所有其他服务器上。
    if len(StaticHosts) > 0 {
        WaitSyncFiles   = append(WaitSyncFiles, uploadRoot+"/"+savePath+"/"+saveName+"."+extension)
    }
    // attach.Self = attach
    // err     = attach.Save(true)
    // if err != nil { return nil, err }

    return attach, nil
}


