package main

import (
	"image"
	"github.com/phachon/graphics-go/graphics"
	"os"
	"image/png"
	"image/jpeg"
	"errors"
)

type Image struct {

}

func NewImage() *Image {
	return &Image{}
}

// 按宽度和高度进行比例缩放
func (img *Image) Scaling(sourceImage string, saveImage string, width int, height int) (err error) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	
	src, filetype, err := LoadImage(sourceImage)
	if err != nil {
		return
	}
	err = graphics.Thumbnail(dst, src)
	if err != nil {
		return
	}
	err = SaveImage(saveImage, dst, filetype)
	return
}

// 根据文件名打开图片,并编码,返回编码对象和文件类型
// Load a image by a filename and return it's type,such as png
func LoadImage(path string) (img image.Image, filetype string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	img, filetype, err = image.Decode(file)
	return
}

// 将编码对象存入文件中
// save a image object into a file just support png and jpg
func SaveImage(path string, img *image.RGBA, filetype string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	if filetype == "png" {
		err = png.Encode(file, img)
	} else if filetype == "jpeg" {
		err = jpeg.Encode(file, img, nil)
	} else {
		err = errors.New("ext of filename not support")
	}
	defer file.Close()
	return
}
