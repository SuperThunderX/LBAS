package main

import (
	"fmt"
	"image"
	"testing"
	"time"

	"github.com/digisan/gotk"
)

func TestDraw(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	img := loadImg("./in/mel.png")

	imagerect := img.Bounds()
	fmt.Println(img.Bounds())

	left, top, right, bottom := imagerect.Min.X, imagerect.Min.Y, imagerect.Max.X, imagerect.Max.Y

	roi := ROIgray(img, left, top, right/2, bottom/2)

	roiClr := ROIrgba(img, left, top, right/2, bottom/2)
	savePNG(roiClr, "./out/color.png")

	roiCMYK := ROIcmyk(img, left, top, right/2, bottom/2)
	savePNG(roiCMYK, "./out/cmyk.png")

	roi1 := DrawHLine(roi, 30, 50, 500, nil)
	roi1 = DrawVLine(roi1, 40, 50, 1500, nil)
	roi1 = DrawRect(roi1, 60, 60, 300, 300, nil)

	savePNG(roi1, "./out/line.png")

	m, _, _ := histogram(roi1.Pix)
	hImg := DrawHisto(m, Peaks(m, 3, 1, -1), Bottoms(m, 3, 1, -1))
	savePNG(hImg, "./out/histo.png")
}

func TestDrawSpline(t *testing.T) {
	img := loadImg("./in/mel.png")
	img = DrawSpline(img,
		[]image.Point{{200, 600}, {100, 300}, {300, 400}, {400, 700}},
		5,
		"R",
		"",
	)
	savePNG(img, "./out/spline-image.png")
}
