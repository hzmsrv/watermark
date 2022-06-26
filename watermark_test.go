package watermark

import "testing"

func Test1(t *testing.T) {
	img, err := New(
		// NOTE 添加修正DEBUG编译文件路径失败！
		WithMark(&mark{
			File:     "./watermark.png",
			Opacity:  30,
			SizeRate: 30,
			Position: "Center",
			Image:    nil,
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	img.MarkFile("test.jpg")
}
