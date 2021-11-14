package main

import (
	// _ "image/jpg"
	"fmt"
	_ "image/png"
	"testing"
	"time"

	"github.com/digisan/gotk"
)

func TestSplitRGBA(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	img := loadImg("./in/sample/1.jpg")
	fmt.Println(img.Bounds())

	r, g, b, a := SplitRGBA(img)
	saveJPG(r, "./out/r.jpg")
	saveJPG(g, "./out/g.jpg")
	saveJPG(b, "./out/b.jpg")
	saveJPG(a, "./out/a.jpg")

	// ///

	com := CompositeRGBA(r, g, b, a)
	saveJPG(com, "./out/com1.png")
}
