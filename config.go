package watermark

import (
	"errors"
	"image"
	"os"
	"path/filepath"

	"github.com/volts-dev/volts/config"
	"github.com/volts-dev/volts/logger"
)

/*
// 水印的位置
const (
	TopLeft Pos = iota
	TopMiddle
	TopRight
	BottomLeft
	BottomMiddle
	BottomRight
	RightMiddle
	LeftMiddle
	Center
)
*/
type (
	//Pos    int
	Option func(*Config)
	Config struct {
		*config.Config `mapstructure:"-"`
		AppPath        string
		OutputExt      string // 输出文件类型
		Height         int
		Width          int
		Marks          []*mark
	}

	mark struct {
		File     string
		Opacity  int
		SizeRate int
		Position string
		Image    image.Image `mapstructure:"-"` // 水印图片
	}
)

var (
	// ErrUnsupportedWatermarkType 不支持的水印类型
	ErrUnsupportedWatermarkType = errors.New("不支持的水印类型")

	log = logger.New("watermark")
)

/*
func (self Pos) String() string {
	switch self {
	case TopLeft:
		return "TopLeft"
	case TopMiddle:
		return "TopMiddle"
	case TopRight:
		return "TopRight"
	case BottomLeft:
		return "BottomLeft"
	case BottomMiddle:
		return "BottomMiddle"
	case BottomRight:
		return "BottomRight"
	case RightMiddle:
		return "RightMiddle"
	case LeftMiddle:
		return "LeftMiddle"
	case Center:
		return "Center"
	default:
		return ""
	}
}
*/
func newConfig(opts ...Option) *Config {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	cfg := &Config{
		Config:  config.New(),
		AppPath: exPath,
		Height:  4, // 按比例
		Width:   4,
	}
	cfg.Init(opts...)
	err = cfg.Load(filepath.Join(exPath, config.CONFIG_FILE_NAME))
	if err != nil {
		log.Err(err)
	}

	err = cfg.UnmarshalField("watermark", &cfg)
	if err != nil {
		log.Err(err)
	}

	//cfg.UnmarshalField("marks", &cfg.Marks)
	if len(cfg.Marks) == 0 {
		cfg.Marks = []*mark{{
			File:     filepath.Join(exPath, "/watermark.png"),
			Opacity:  30,
			SizeRate: 30,
			Position: "Center",
			Image:    nil,
		}}

		cfg.SetValue("watermark.marks", cfg.Marks)
		cfg.Save()
	}

	return cfg
}

func (self *Config) Init(opts ...Option) {
	for _, opt := range opts {
		opt(self)
	}
}

func WithMark(m *mark) Option {
	return func(cfg *Config) {
		cfg.Marks = append(cfg.Marks, m)
	}
}

func WithSize(height, width int) Option {
	return func(cfg *Config) {
		cfg.Height = height
		cfg.Width = width
	}
}
