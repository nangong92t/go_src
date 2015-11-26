package modules

import (
    "errors"
    "strings"
    "time"
    "fmt"
    "os"
    "io/ioutil"
    "crypto/md5"
    "encoding/hex"
    "path/filepath"
    "github.com/coopernurse/gorp"
    "github.com/revel/revel"
    "std/data-service/std/app/models"
)

type Attach struct {
    Db *gorp.DbMap
}

func NewAttach(db *gorp.DbMap) *Attach {
    return &Attach{
        Db: db,
    }
}

// 执行图片上传
func (m *Attach) Upload(req *revel.Request, fileName string) (*models.Attach, error) {
    // 获取配置文件中上传根路径
    revel.Config.SetSection( revel.RunMode )
    uploadRoot, found   := revel.Config.String("uploadpath")

    if !found {
        return nil, errors.New("Sorry, no found out the upload root path config")
    }

    // 获取上传文件数据
    file, handler, err  := req.FormFile(fileName)
    if err != nil {
        return nil, err
    }

    // 获取文件扩展名
    extension           := strings.ToLower( filepath.Ext(handler.Filename) )
    // 获取文件的存储路径
    savePath, saveName  := m.MakeFolderPath(uploadRoot, handler.Filename)

    // 读取图片二进制文件.
    data, err2          := ioutil.ReadAll(file)
    if err2 != nil {
        return nil, err
    }

    // 获取图片尺寸
    size                := int64(len(data))

    // 写入文件
    err = ioutil.WriteFile(uploadRoot+"/"+savePath+"/"+saveName+extension, data, 0644)
    if err != nil {
        return nil, err
    }

    attach  := &models.Attach {
        Size:   size,
        Extention: extension,
        Name:   handler.Filename,
        Md5:   "",
        SaveHost: "",
        SavePath: savePath,
        SaveName: saveName,
        Created: time.Now().Unix(),
    }

    err = m.Db.Insert(attach)

    if err != nil {
        return nil, err
    }

    return attach, nil
}

// 根据文件名生成文件路径, 如果文件路径不存在并创建
func (m *Attach) MakeFolderPath(rootPath string, fileName string) (string, string) {
    secs            := time.Now().Unix()

    h               := md5.New()
    h.Write([]byte(fileName + fmt.Sprintf("%s", secs)))
    fileName        = fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
    pathstrs        := []rune(fileName)
    f1, f2, f3      := string(pathstrs[0]), string(pathstrs[1:4]), string(pathstrs[5:10])

    tPath           := rootPath + "/" + f1
    isExists, _     := m.IsFolder(tPath)
    if !isExists {
        os.Mkdir(tPath, 0755)
    }

    tPath           = tPath + "/" + f2
    isExists, _     = m.IsFolder(tPath)
    if !isExists {
        os.Mkdir(tPath, 0755)
    }

    tPath           = tPath + "/" + f3
    isExists, _     = m.IsFolder(tPath)
    if !isExists {
        os.Mkdir(tPath, 0755)
    }

    path            := f1 + "/" + f2 + "/" + f3

    return path, fileName
}

// 检测文件夹是否存在.
func (m *Attach) IsFolder(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

// 获取图片url地址
func (m *Attach) GetAttachUrl(attach *models.Attach) (string, error) {
    attachUrl   := ""

    if attach.Id == 0 {
        return "", nil
    }

    // 获取配置文件中静态文件域名.
    revel.Config.SetSection( revel.RunMode )
    staticDomain, found   := revel.Config.String("staticdomain")

    if !found {
        return "", errors.New("Sorry, no found out the static domain config")
    }

    attachUrl   =   staticDomain +
                    "/uploads/" + attach.SavePath +
                    "/" + attach.SaveName +
                    attach.Extention

    return attachUrl, nil
}

// 通过Attach Id获取对应Map数据.
func (m *Attach) GetAttachMapById(attachIds []string) (map[int64]models.Attach, error) {
    idLen   := len(attachIds)
    urlMap  := make(map[int64]models.Attach, idLen)

    var attachs []models.Attach

    _, err  := m.Db.Select(&attachs, "select * from attach where attach_id in (" + strings.Join(attachIds, ", ") + ")")

    if err != nil { return nil, err }

    for _, one := range attachs {
        one.Uri, _ = m.GetAttachUrl(&one)
        urlMap[ one.Id ]    = one
    }

    return urlMap, nil
}

// 获取一个attach数据.
func (m *Attach) GetOneAttach(id int64) *models.Attach {
    attach := &models.Attach{}

    if id == 0 { return attach }

    m.Db.SelectOne(attach, "select * from attach where attach_id=? limit 1", id)
    return attach
}

