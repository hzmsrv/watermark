package watermark

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
)

type (
	output struct {
		buf *image.NRGBA64
	}
)

func (self *output) AsBytes() []byte {
	return nil
}

func (self *output) AsJpg() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := jpeg.Encode(buf, self.buf, &jpeg.Options{Quality: jpeg.DefaultQuality})
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (self *output) AsPng() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := png.Encode(buf, self.buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (self *output) AsWebp() (*bytes.Buffer, error) {
	buf := bytes.NewBuffer([]byte{})
	err := webp.Encode(buf, self.buf, &webp.Options{Lossless: true})
	if err != nil {
		return nil, err
	}

	return buf, nil
}
