package img

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"image"
	"os"
	"path"
)

// 获取图片方向信息
func Orientation(filename string) (int64, error) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return 0, err
	}

	// 解析器注册
	exif.RegisterParsers(mknote.All...)

	// 解析图片
	x, err := exif.Decode(f)
	if x != nil {
		// 获取图片方向
		t, _ := x.Get(exif.Orientation)
		return t.Int64(0)
	}

	return 0, nil
}

// 图片旋转
func RotateImage(src string) error {
	// 获取图片方向信息
	o, err := Orientation(src)
	if err != nil {
		logs.Error(fmt.Errorf("获取图片%s方向信息失败 %s", src, err))
		return fmt.Errorf("获取图片%s方向信息失败 %s", src, err)
	}
	return RotateImageByOrientation(src, o)
}

// 图像方向信息orientation
// 相机常用翻转情况有4种
// 1: 头朝上
// 6: 头朝右
// 3: 头朝下
// 8: 头朝左
// 参考: https://www.impulseadventure.com/photo/exif-orientation.html
// 图片旋转(指定方向)
func RotateImageByOrientation(src string, orientation int64) error {
	openImage, err := imaging.Open(src)
	if err != nil {
		logs.Error(fmt.Errorf("读取图片%s失败 %s", src, err))
		return fmt.Errorf("读取图片%s失败 %s", src, err)
	}
	saveTmp := src + ".RotateImage" + path.Ext(src)
	var rotateImage *image.NRGBA
	switch orientation {
	case 1:
		// Do nothing
		break
	case 2:
		// 水平翻转(不常用)
		rotateImage = imaging.FlipH(openImage)
		break
	case 3:
		// 逆时针旋转180度
		rotateImage = imaging.Rotate180(openImage)
		break
	case 4:
		// 垂直翻转
		rotateImage = imaging.FlipV(openImage)
		break
	case 5:
		// 垂直翻转并且逆时针旋转90度
		rotateImage = imaging.Transverse(openImage)
		break
	case 6:
		// 逆时针旋转270度
		rotateImage = imaging.Rotate270(openImage)
		break
	case 7:
		// 水平翻转并且逆时针选择90度
		rotateImage = imaging.Transpose(openImage)
		break
	case 8:
		// 逆时针旋转90度
		rotateImage = imaging.Rotate90(openImage)
		break
	}
	if rotateImage != nil {
		// 保存图片到临时图片
		err = imaging.Save(rotateImage, saveTmp)
		if err != nil {
			logs.Error(fmt.Errorf("保存到临时图片%s失败 %s", saveTmp, err))
			return fmt.Errorf("保存到临时图片%s失败 %s", saveTmp, err)
		}
		// 重命名临时图片
		err := os.Rename(saveTmp, src)
		if err != nil {
			logs.Error(fmt.Errorf("图片重命名%s => %s失败 %s", saveTmp, src, err))
			return fmt.Errorf("图片重命名%s => %s失败 %s", saveTmp, src, err)
		}
	} else {
		logs.Debug("图片无须旋转")
	}
	return nil
}
