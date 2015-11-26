// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

// 转用于公共服务接口.

package controllers

import (
    r "github.com/revel/revel"
    "maintance/staticserver/app/components/uploader"
)

//
// for Public.
//
type Public struct {
    AbstractController
}

// 根据需求实时获取缩略图片
// way string 缩略图生成方式:
//      'auto': 原图全显优先，获取width或height最大，原图不会被截断
//      'full': 图片填充优先，将widthd和height，按照比例全部填充， 原图将会被截断
// 访问形式如: /thumb/full/180/120/0/b6e/cf2fa/0b6eacf2fa64cbdc21cceda145a56b53
// 或者: /1/public/thumb?way=auto&width=150&height=200&picturePath=0/b6e/cf2fa/0b6eacf2fa64cbdc21cceda145a56b53.jpg
// 
func (c *Public) Thumb(way string, width uint, height uint, picturePath string) r.Result {
    pictureBytes, err := uploader.Thumb(way, width, height, picturePath)
    if err != nil { return c.Render(nil, err) }

    return c.RenderFile(pictureBytes, r.Inline)
}

// for upload the picture.
//  just support: multipart/form-data
func (c *Public) Upload() r.Result {
    attach, err := uploader.Upload(c.Request.Request)
    return c.Render(attach, err)
}

// 自动同步新上传静态文件到多服务器。
func (c *Public) UploadSync() r.Result {
    err := uploader.UploadSync(c.Request.Request)
    return c.Render(nil, err)
}
