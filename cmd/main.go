package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"volts-dev/watermark"

	"github.com/volts-dev/volts/logger"
)

var log = logger.New("watermark")

func main() {
	log.Info("明赛科技")

	if time.Now().Year() >= 2024 {
		return
	}

	// First element in os.Args is always the program name,
	// So we need at least 2 arguments to have a file name argument.
	if len(os.Args) < 2 {
		log.Fatal("Missing parameter, provide file name!")
		return
	}
	file := os.Args[1]
	fp := filepath.Dir(file)
	fileName := filepath.Base(file)

	img, err := watermark.New()
	if err != nil {
		log.Panicf("添加水印失败！&v", err)
	}

	data, err := img.MarkFile(file)
	if err != nil {
		log.Panicf("添加水印失败！&v", err)
	}

	buf, err := data.AsJpg()
	if err != nil {
		log.Panicf("添加水印失败！&v", err)
	}

	//ext := filepath.Ext(file)
	//if ext == ".webp" {
	//	ext = ".jpg"
	//}
	fp = filepath.Join(fp, fmt.Sprintf("%s_%s%s", strings.Split(fileName, ".")[0], "marked", ".jpg"))
	imgw, _ := os.Create(fp)
	defer imgw.Close()
	imgw.ReadFrom(buf)
}
