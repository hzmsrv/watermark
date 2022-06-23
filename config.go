package watermark

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/volts-dev/volts/logger"
)

type (
	Option func(*Config)
	Config struct {
		AppPath       string
		WaterMarkFile string
		Output        string
		Opacity       int
		SizeRate      int
		Height        int
		Width         int
	}
)

var (
	// ErrUnsupportedWatermarkType 不支持的水印类型
	ErrUnsupportedWatermarkType = errors.New("不支持的水印类型")
	// ErrWatermarkTooLarge 当水印位置距离右下角的范围小于水印图片时，返回错误。
	ErrWatermarkTooLarge = errors.New("水印太大")

	log = logger.New("watermark")
)

func newConfig(opts ...Option) *Config {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	cfg := &Config{
		AppPath:       exPath,
		WaterMarkFile: filepath.Join(exPath, "/watermark.png"),
		Height:        1080,
		Width:         1080,
		SizeRate:      30, // 30%
		Opacity:       30,
	}
	cfg.Init(opts...)
	return cfg
}

func (self *Config) Init(opts ...Option) {
	for _, opt := range opts {
		opt(self)
	}
}

func WithPosition() Option {
	return func(cfg *Config) {

	}
}

func WithSize(height, width int) Option {
	return func(cfg *Config) {
		cfg.Height = height
		cfg.Width = width
	}
}
