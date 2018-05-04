package utils

import (
	"os"
	"image"
	"image/png"
	"image/jpeg"
	"errors"
	"github.com/phachon/graphics-go/graphics"
)

type Imager struct {

}

func NewImager() *Imager {
	return &Imager{}
}

// 按宽度和高度进行比例缩放
func (imager *Imager) Scaling(sourceImage string, saveImage string, width int, height int) (err error) {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	src, fileType, err := imager.Decode(sourceImage)
	if err != nil {
		return
	}
	err = graphics.Thumbnail(dst, src)
	if err != nil {
		return
	}
	err = imager.SaveImage(saveImage, dst, fileType)
	return
}

// 对图片解码
func (imager *Imager) Decode(imagePath string) (img image.Image, fileType string, err error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return
	}
	defer file.Close()
	img, fileType, err = image.Decode(file)
	return
}

// 保存一个 image 对象为文件
func (imager *Imager) SaveImage(savePath string, img *image.RGBA, fileType string) (err error) {
	file, err := os.Create(savePath)
	if err != nil {
		return
	}
	defer file.Close()

	if fileType == "png" {
		err = png.Encode(file, img)
	} else if fileType == "jpeg" {
		err = jpeg.Encode(file, img, nil)
	} else {
		err = errors.New("ext of filename not support")
	}

	return
}