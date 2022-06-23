package watermark

import "testing"

func Test1(t *testing.T) {
	img, err := New()
	if err != nil {
		t.Fatal(err)
	}

	img.MarkFile("test.jpg")
}
