package watermark

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	//"golang.org/x/image/webp"
)

// 水印的位置
const (
	TopLeft Pos = iota
	TopMiddle
	TopRight
	BottomLeft
	BottomMiddle
	BottomRight
	Center
)

type (
	Pos       int
	watermark struct {
		Config   *Config
		image    image.Image // 水印图片
		Position Pos
		padding  int
	}
)

func New(opts ...Option) (*watermark, error) {
	cfg := newConfig(opts...)

	f, err := os.Open(cfg.WaterMarkFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var img image.Image
	var gifImg *gif.GIF
	switch strings.ToLower(filepath.Ext(cfg.WaterMarkFile)) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	case ".gif":
		gifImg, err = gif.DecodeAll(f)
		img = gifImg.Image[0]
	default:
		return nil, ErrUnsupportedWatermarkType
	}

	if err != nil {
		return nil, err
	}

	wm := &watermark{
		Config:   cfg,
		image:    img,
		Position: Center,
	}
	/*
		imgb, _ := os.Open("image.jpg")
		img, _ := jpeg.Decode(imgb)
		defer imgb.Close()

		wmb, _ := os.Open("watermark.png")
		watermark, _ := png.Decode(wmb)
		defer wmb.Close()

		offset := image.Pt(200, 200)
		b := img.Bounds()
		m := image.NewRGBA(b)
		draw.Draw(m, b, img, image.ZP, draw.Src)
		draw.Draw(m, watermark.Bounds().Add(offset), watermark, image.ZP, draw.Over)

		imgw, _ := os.Create("watermarked.jpg")
		jpeg.Encode(imgw, m, &jpeg.Options{jpeg.DefaultQuality})
		defer imgw.Close()
	*/
	return wm, nil
}

func (self *watermark) Mark(src io.ReadWriteSeeker, ext string, opts ...Option) (out *output, err error) {
	var srcImg image.Image

	ext = strings.ToLower(ext)
	switch ext {
	case ".webp":
		//	x, x1 := webp.Decode(src)
		srcImg, err = webp.Decode(src) // image.Image(x), x1
	case ".gif":
		//return w.markGIF(src) // GIF 另外单独处理
	case ".jpg", ".jpeg":
		srcImg, err = jpeg.Decode(src)
	case ".png":
		srcImg, err = png.Decode(src)
	default:
		return nil, ErrUnsupportedWatermarkType
	}
	if err != nil {
		return nil, err
	}

	// 切割大小为正方形
	srcImg, err = self.resize(srcImg)
	if err != nil {
		return nil, err
	}

	// 伸缩水印大小
	var rate float64 = 1.0
	//bigint := self.image.Bounds().Max.X
	if self.image.Bounds().Max.X < self.image.Bounds().Max.Y {
		//	bigint = self.image.Bounds().Max.Y
		rate = float64(self.image.Bounds().Dy()) / (float64(srcImg.Bounds().Dy()) * (float64(self.Config.SizeRate) * 0.01))
	} else if self.image.Bounds().Max.X > self.image.Bounds().Max.Y {
		rate = float64(self.image.Bounds().Dx()) / (float64(srcImg.Bounds().Dx()) * (float64(self.Config.SizeRate) * 0.01))

	}

	heigh := uint(math.Ceil(float64(self.image.Bounds().Dy()) / rate))
	width := uint(math.Ceil(float64(self.image.Bounds().Dx()) / rate))
	mark := resize.Resize(width, heigh, self.image, resize.Lanczos3)
	point := self.getPoint(srcImg, mark)

	//	if err = self.checkTooLarge(point, bound); err != nil {
	//		return err
	//	}

	dstImg := image.NewNRGBA64(srcImg.Bounds())
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Point{}, draw.Over)
	draw.DrawMask(dstImg, dstImg.Bounds(), mark, point, image.NewUniform(color.Alpha{uint8(float32(self.Config.Opacity) * 2.55)}), point, draw.Over)
	if _, err = src.Seek(0, 0); err != nil {
		return nil, err
	}

	return &output{
		buf: dstImg,
	}, nil

	/*
		imgw, _ := os.Create(fmt.Sprintf("watermarked_%s.jpg", "marked"))
		defer imgw.Close()
		switch ext {
		case ".webp":
			//return webp.Encode(imgw, dstImg, &webp.Options{Lossless: true})
			return jpeg.Encode(imgw, dstImg, &jpeg.Options{jpeg.DefaultQuality})

		case ".jpg", ".jpeg":
			//	 jpeg.Encode(src, dstImg, nil)
			return jpeg.Encode(imgw, dstImg, &jpeg.Options{jpeg.DefaultQuality})
		case ".png":
			return png.Encode(imgw, dstImg)
		default:
			return ErrUnsupportedWatermarkType
		}
	*/
}

// MarkFile 给指定的文件打上水印
func (self *watermark) MarkFile(path string, opts ...Option) (*output, error) {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return self.Mark(file, strings.ToLower(filepath.Ext(path)), opts...)
}

func (self *watermark) checkTooLarge(start image.Point, dst image.Rectangle) error {
	// 允许的最大高宽
	width := dst.Dx() - start.X - self.padding
	height := dst.Dy() - start.Y - self.padding

	if width < self.image.Bounds().Dx() || height < self.image.Bounds().Dy() {
		return ErrWatermarkTooLarge
	}
	return nil
}

func (self *watermark) resize(img image.Image) (image.Image, error) {
	var heightwidth uint
	if img.Bounds().Dx() > img.Bounds().Dy() {
		// Set the expected size that you want:
		heightwidth = uint(img.Bounds().Dy())
	} else if img.Bounds().Dx() < img.Bounds().Dy() {
		heightwidth = uint(img.Bounds().Dx())
	}

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	dst := resize.Resize(heightwidth, heightwidth, img, resize.Lanczos3)

	return dst, nil
}

func (self *watermark) getPoint(src, mark image.Image) image.Point {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	var point image.Point
	switch self.Position {
	case TopLeft:
		point = image.Point{X: -self.padding, Y: -self.padding}
	case TopMiddle:
		point = image.Point{
			X: -self.padding,
			Y: -(height - self.padding - mark.Bounds().Dy()) / 2,
		}
	case TopRight:
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()),
			Y: -self.padding,
		}
	case BottomLeft:
		point = image.Point{
			X: -self.padding,
			Y: -(height - self.padding - mark.Bounds().Dy()),
		}
	case BottomMiddle:
		point = image.Point{
			X: -self.padding,
			Y: -(height - self.padding - mark.Bounds().Dy()) / 2,
		}
	case BottomRight:
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()),
			Y: -(height - self.padding - mark.Bounds().Dy()),
		}
	case Center:
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()) / 2,
			Y: -(height - self.padding - mark.Bounds().Dy()) / 2,
		}
	default:
		panic("无效的 pos 值")
	}

	return point
}
