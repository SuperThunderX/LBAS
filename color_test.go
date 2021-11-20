package main

import (
	// _ "image/jpg"

	"fmt"
	_ "image/png"
	"os"
	"testing"
	"time"

	"github.com/digisan/gotk"
)

func TestSplitRGBA(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	os.MkdirAll("./out", os.ModePerm)

	img := loadImage("./out/sample-roi.jpg")
	fmt.Println(img.Bounds())

	r, g, b, a := SplitRGBA(img)
	saveJPG(r, "./out/r.jpg")
	saveJPG(g, "./out/g.jpg")
	saveJPG(b, "./out/b.jpg")
	saveJPG(a, "./out/a.jpg")

	// ///

	com := CompositeRGBA(r, g, b, a)
	saveJPG(com, "./out/com1.jpg")
}
