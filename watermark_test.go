package watermark

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test1(t *testing.T) {
	img, err := New(
		// NOTE 添加修正DEBUG编译文件路径失败！
		WithMark(mark{
			File:     "./frame.png",
			Opacity:  30,
			SizeRate: 30,
			Position: "Center",
			Image:    nil,
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	data, err := img.MarkFile("test.webp")
	if err != nil {
		t.Fatal(err)
	}

	buf, err := data.AsJpg()
	if err != nil {
		log.Panicf("添加水印失败！&v", err)
	}

	//ext := filepath.Ext(file)
	//if ext == ".webp" {
	//	ext = ".jpg"
	//}
	fp := filepath.Join("./", fmt.Sprintf("%s_%s%s", strings.Split("test", ".")[0], "marked", ".jpg"))
	//fmt.Println(fp)

	imgw, _ := os.Create(fp)
	defer imgw.Close()
	imgw.ReadFrom(buf)
}
