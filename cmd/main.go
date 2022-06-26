package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"volts-dev/watermark"

	//"volts-dev/watermark"

	"github.com/volts-dev/volts/logger"
)

var log = logger.New("watermark")

func main() {
	log.Info("明赛科技")

	if time.Now().Year() == 2023 {
		return
	}
	//logg, _ := os.Create("./log")
	//defer logg.Close()

	// First element in os.Args is always the program name,
	// So we need at least 2 arguments to have a file name argument.
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}
	file := os.Args[1]
	fp := filepath.Dir(file)
	fileName := filepath.Base(file)
	fmt.Println(file)
	fmt.Println(fp, " ", fileName)
	//logg.WriteString(file)
	//	logg.WriteString(fp)
	//logg.WriteString(fileName)

	img, err := watermark.New()
	if err != nil {
		log.Panicf("添加水印失败！&v", err)
		//logg.WriteString(err.Error())
	}

	data, err := img.MarkFile(file)
	if err != nil {
		log.Panicf("添加水印失败！&v", err)
		//logg.WriteString(err.Error())
	}

	buf, err := data.AsJpg()
	if err != nil {
		log.Panicf("添加水印失败！&v", err)
		//logg.WriteString(err.Error())
	}

	fp = filepath.Join(fp, fmt.Sprintf("%s_%s%s", strings.Split(fileName, ".")[0], "marked", filepath.Ext(file)))
	fmt.Println(fp)

	imgw, _ := os.Create(fp)
	defer imgw.Close()
	imgw.ReadFrom(buf)
}
