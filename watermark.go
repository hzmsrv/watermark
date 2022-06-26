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
	"github.com/oliamb/cutter"
	//"golang.org/x/image/webp"
)

type (
	watermark struct {
		Config *Config
		//image    image.Image // 水印图片

		padding int
	}
)

func New(opts ...Option) (*watermark, error) {
	cfg := newConfig(opts...)

	for _, mk := range cfg.Marks {
		f, err := os.Open(mk.File)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		var img image.Image // 水印图片
		var gifImg *gif.GIF
		switch strings.ToLower(filepath.Ext(mk.File)) {
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

		mk.Image = img
	}

	wm := &watermark{
		Config: cfg,
		//Position: Center,
	}

	return wm, nil
}

func (self *watermark) Mark(src io.ReadWriteSeeker, ext string, opts ...Option) (out *output, err error) {
	var srcImg image.Image

	ext = strings.ToLower(ext)
	switch ext {
	case ".webp":
		srcImg, err = webp.Decode(src)
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
	dstImg := image.NewNRGBA64(srcImg.Bounds())
	draw.Draw(dstImg, dstImg.Rect, srcImg, dstImg.Rect.Min, draw.Src) // 画布
	for _, mk := range self.Config.Marks {
		// 伸缩水印大小
		var scalerate float64 = 1.0
		if mk.Image.Bounds().Max.X < mk.Image.Bounds().Max.Y {
			scalerate = float64(mk.Image.Bounds().Dy()) / (float64(srcImg.Bounds().Dy()) * (float64(mk.SizeRate) * 0.01))
		} else if mk.Image.Bounds().Max.X > mk.Image.Bounds().Max.Y {
			scalerate = float64(mk.Image.Bounds().Dx()) / (float64(srcImg.Bounds().Dx()) * (float64(mk.SizeRate) * 0.01))

		}

		heigh := uint(math.Ceil(float64(mk.Image.Bounds().Dy()) / scalerate))
		width := uint(math.Ceil(float64(mk.Image.Bounds().Dx()) / scalerate))
		mark := resize.Resize(width, heigh, mk.Image, resize.Lanczos3)
		point := self.getPoint(srcImg, mark, mk.Position)

		draw.DrawMask(dstImg, dstImg.Bounds(), mark, point, image.NewUniform(color.Alpha{uint8(float32(mk.Opacity) * 2.55)}), point, draw.Over)

	}
	if _, err = src.Seek(0, 0); err != nil {
		return nil, err
	}

	return &output{
		buf: dstImg,
	}, nil
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

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// TODO 实现尺寸比例
func (self *watermark) resize(img image.Image) (image.Image, error) {
	var miniEdge int
	/*var Width, Height int
	if self.Config.Height == self.Config.Width { // 如果是正方形1:1比例切割

		if img.Bounds().Dx() > img.Bounds().Dy() {
			// Set the expected size that you want:
			miniEdge = img.Bounds().Dy()
		} else {
			miniEdge = img.Bounds().Dx()
		}
		Width, Height = miniEdge, miniEdge
	} else {
		// 如果要求比例高大于宽
		if self.Config.Height > self.Config.Width {
			rate := (img.Bounds().Dy() / self.Config.Height)
			if img.Bounds().Dx() > img.Bounds().Dy() {
				// 当比例高最大且实图高最小时 高不变计算宽度
				Width = rate * self.Config.Height
				Height = rate * self.Config.Width
			} else if img.Bounds().Dx() < img.Bounds().Dy() {
				Height = rate * self.Config.Height
				Width = rate * self.Config.Width
			}
		} else {
			rate := (img.Bounds().Dx() / self.Config.Height)
			if img.Bounds().Dx() > img.Bounds().Dy() {
				// 当比例高最大且实图高最小时 高不变计算宽度
				Width = rate * self.Config.Height
				Height = rate * self.Config.Width
			} else if img.Bounds().Dx() < img.Bounds().Dy() {
				Height = rate * self.Config.Height
				Width = rate * self.Config.Width
			}
		}

	}
	*/
	if img.Bounds().Dx() > img.Bounds().Dy() {
		// Set the expected size that you want:
		miniEdge = img.Bounds().Dy()
	} else if img.Bounds().Dx() < img.Bounds().Dy() {
		miniEdge = img.Bounds().Dx()
	}
	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	//dst := resize.Resize(heightwidth, heightwidth, img, resize.Lanczos3)

	return cutter.Crop(img, cutter.Config{
		Width:   miniEdge,
		Height:  miniEdge,
		Mode:    cutter.Centered,
		Options: cutter.Ratio & cutter.Copy, // Copy is useless here
	})
}

func (self *watermark) getPoint(src, mark image.Image, pos string) image.Point {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()
	var point image.Point
	switch strings.ToLower(pos) {
	case "top_left":
		point = image.Point{X: -self.padding, Y: -self.padding}
	case "top_middle":
		point = image.Point{
			//X: -self.padding,
			//Y: -(height - self.padding - mark.Bounds().Dy()) / 2,
			X: -(width - self.padding - mark.Bounds().Dx()) / 2,
			Y: -self.padding,
		}
	case "top_right":
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()),
			Y: -self.padding,
		}
	case "bottom_left":
		point = image.Point{
			X: -self.padding,
			Y: -(height - self.padding - mark.Bounds().Dy()),
		}
	case "bottom_middle":
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()) / 2,
			Y: -(height - self.padding - mark.Bounds().Dy()),
		}
	case "bottom_right":
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()),
			Y: -(height - self.padding - mark.Bounds().Dy()),
		}
	case "right_middle":
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()),
			Y: -(height - self.padding - mark.Bounds().Dy()) / 2,
		}
	case "left_middle":
		point = image.Point{
			X: -self.padding,
			Y: -(height - self.padding - mark.Bounds().Dy()) / 2,
		}
	case "center":
		point = image.Point{
			X: -(width - self.padding - mark.Bounds().Dx()) / 2,
			Y: -(height - self.padding - mark.Bounds().Dy()) / 2,
		}
	default:
		panic("无效的 pos 值")
	}

	return point
}
