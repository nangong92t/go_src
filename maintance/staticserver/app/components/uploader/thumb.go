// Copyright 2014, SuccessfumMatch Teq Inc. All rights reserved.
// Author TonyXu<tonycbcd@gmail.com>,
// Build on dev-0.0.1
// MIT Licensed

package uploader

import (
	"code.google.com/p/graphics-go/graphics"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/nfnt/resize"
	"github.com/revel/revel"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"syscall"
)

// 根据需求实时获取缩略图片
// way string 缩略图生成方式:
//      'auto': 原图全显优先，获取width或height最大，原图不会被截断
//      'full': 图片填充优先，将widthd和height，按照比例全部填充， 原图将会被截断
//
func Thumb(way string, width uint, height uint, picturePath string) (*os.File, error) {
	/* 获取配置文件中上传根路径 */
	revel.Config.SetSection(revel.RunMode)
	uploadRoot, _ := revel.Config.String("uploadpath")

	ext := filepath.Ext(picturePath)

	/* 检测是否是获取一个远端url图片的缩率图 */
	picUrl, err := url.ParseRequestURI(picturePath)
	isRemotePic := false
	//if err == nil && picUrl.Scheme == "http" && (picUrl.Host == "img.show.wepiao.com" || picUrl.Host == "img.wxmovie.com" || picUrl.Host == "static.show.wepiao.com" || picUrl.Host == "static.task.18.tl:86") {
	if err == nil && picUrl.Scheme == "http" {
		isRemotePic = true
		uploadRoot += "/upload"
		savePath, saveName := HelperForFilePathCreator(uploadRoot, picturePath, false)
		picturePath = savePath + "/" + saveName + ext
	}

	if way != "full" {
		way = "auto"
	}

	fname := uploadRoot + "/" + picturePath
	thumbFile := fmt.Sprintf("%s_%s_%v_%v%v", fname, way, width, height, ext)

	/* 检查有没有现成的缩略图 */
	_, err = os.Stat(thumbFile)
	if err == nil { // 如果有，直接读出
		thumb, _ := os.Open(thumbFile)
		return thumb, nil
	}
	// fmt.Printf("没有现成缩率图\n")

	/* 没有现成的缩略图，需要实施生成, 查找原图片文件是否存在 */
	finfo, err := os.Stat(fname)
	if err != nil {
		// 如果是获取远端图片，那就先获取.
		if isRemotePic {
			// fmt.Printf("是远端图片\n")

			client := http.Client{}
			reqImg, err := client.Get(picUrl.Scheme + "://" + picUrl.Host + picUrl.Path)
			defer reqImg.Body.Close()
			if err != nil {
				return nil, err
			}
			//如果是4开头的错误状态 则生成默认图片
			if '4' == reqImg.Status[0] {
				fname = "/opt/www/static/pc/img/logo.png"
				way = "auto"
			} else {
				out, err := os.Create(fname)
				if err != nil {
					return nil, err
				}
				defer out.Close()
				io.Copy(out, reqImg.Body)
			}
			finfo, err = os.Stat(fname)
			// fmt.Printf("存储远端图片到: %#v\n", fname)
		} else if os.IsNotExist(err) || err.(*os.PathError).Err == syscall.ENOTDIR {
			return nil, errors.New("20108")
		} else {
			return nil, err
		}
	}
	if finfo.Mode().IsDir() {
		return nil, errors.New("10016")
	}

	file, err := os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("20108")
		}
		return nil, err
	}
	defer file.Close()

	/* 生成缩略图 */
	var img image.Image
	img, _, err = image.Decode(file)

	if err != nil {
		return nil, err
	}

	// 产生缩略图,等比例缩放
	out, errF := os.Create(thumbFile)
	if errF != nil {
		return nil, errors.New("20108")
	}
	defer out.Close()

	var m image.Image
	// fmt.Printf("生成缩率图\n")

	origBounds := img.Bounds()
	origWidth := uint(origBounds.Dx())
	origHeight := uint(origBounds.Dy())
	if way == "full" {
		dst := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
		graphics.Thumbnail(dst, img)
		m = dst
	} else if width >= origWidth && height >= origHeight {
		bg := imaging.New(int(width), int(height), color.NRGBA{255, 255, 255, 255})
		m = imaging.PasteCenter(bg, img)
	} else {
		m = resize.Thumbnail(width, height, img, resize.Lanczos3)
	}

	/* 生成文件 */
	err = jpeg.Encode(out, m, &jpeg.Options{90})
	if err != nil {
		err = png.Encode(out, m)
	}
	if err != nil {
		err = gif.Encode(out, m, nil)
	}
	if err != nil {
		return nil, err
	}

	/* 这里不知道如何将 image.Image类型转换为 io.Writer类型，
	   所以只有暂时在保存后，又重新读了次文件 */
	thumb, _ := os.Open(thumbFile)

	return thumb, nil
}
