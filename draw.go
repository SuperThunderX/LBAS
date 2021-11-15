package main

import (
	"image"
	"image/draw"
	"math"
	"sort"

	"github.com/cnkei/gospline"
	"github.com/digisan/gotk/slice/ti"
	"github.com/digisan/gotk/slice/tu8i"
	"github.com/fogleman/gg"
)

func DrawRect(img *image.Gray, left, top, right, bottom int, paint *image.Gray) *image.Gray {
	if paint == nil {
		paint = image.NewGray(image.Rect(0, 0, right, bottom))
		for i := 0; i < len(paint.Pix); i++ {
			paint.Pix[i] = 0
		}
	}

	draw.Draw(img, image.Rect(left, top, right, top+1), paint, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(left, top, left+1, bottom), paint, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(left, bottom, right, bottom+1), paint, image.Point{0, 0}, draw.Src)
	draw.Draw(img, image.Rect(right, top, right+1, bottom+1), paint, image.Point{0, 0}, draw.Src)
	return img
}

func DrawHLine(img *image.Gray, y, left, right int, paint *image.Gray) *image.Gray {
	if paint == nil {
		paint = image.NewGray(image.Rect(0, 0, right, 1))
		for i := 0; i < len(paint.Pix); i++ {
			paint.Pix[i] = 0
		}
	}

	draw.Draw(img, image.Rect(left, y, right, y+1), paint, image.Point{0, 0}, draw.Src)
	return img
}

func DrawVLine(img *image.Gray, x, top, bottom int, paint *image.Gray) *image.Gray {
	if paint == nil {
		paint = image.NewGray(image.Rect(0, 0, 1, bottom))
		for i := 0; i < len(paint.Pix); i++ {
			paint.Pix[i] = 0
		}
	}

	draw.Draw(img, image.Rect(x, top, x+1, bottom), paint, image.Point{0, 0}, draw.Src)
	return img
}

func DrawHisto(mHisto, mPeak, mBottom map[byte]int) (hImg *image.Gray) {

	_, vs := tu8i.Map2KVs(mHisto, nil, nil)
	maxCnt := ti.Max(vs...)
	r := float64(maxCnt) / float64(255)
	hImg = image.NewGray(image.Rect(0, 0, 256, 256))

	// drawing
	mY := make(map[byte]int)
	for k, v := range mHisto {
		mY[k] = int(float64(v) / r)
	}
	ks, vs := tu8i.Map2KVs(mY, func(i, j byte) bool { return i < j }, nil)
	vs = smooth(vs) // remove noise

	paint := image.NewGray(image.Rect(0, 0, 1, 256))
	for i := 0; i < len(paint.Pix); i++ {
		paint.Pix[i] = 255
	}
	for i := 0; i < len(ks); i++ {
		k, v := ks[i], vs[i]
		DrawVLine(hImg, int(k), 256-v, 256, paint)
	}

	// mark peak
	if len(mPeak) > 0 {
		paintPeak := image.NewGray(image.Rect(0, 0, 1, 50))
		for i := 0; i < len(paintPeak.Pix); i++ {
			paintPeak.Pix[i] = 50
		}
		for x, y := range mPeak {
			y = int(float64(y) / r)
			DrawVLine(hImg, int(x), 256-y+5, 256-y+30, paintPeak)
		}
	}

	// mark bottom
	if len(mBottom) > 0 {
		paintBottom := image.NewGray(image.Rect(0, 0, 1, 50))
		for i := 0; i < len(paintBottom.Pix); i++ {
			paintBottom.Pix[i] = 200
		}
		for x, y := range mBottom {
			y = int(float64(y) / r)
			DrawVLine(hImg, int(x), 256-y-30, 256-y-5, paintBottom)
		}
	}

	return hImg
}

///////////////////////////////////////////////////////////////////////////

func ZipPoints(xs, ys []float64) (pts []image.Point) {
	for i, x := range xs {
		y := ys[i]
		pts = append(pts, image.Point{X: int(x), Y: int(y)})
	}
	return
}

func UnzipPoints(pts []image.Point) (xs, ys []float64) {
	for _, pt := range pts {
		xs = append(xs, float64(pt.X))
		ys = append(ys, float64(pt.Y))
	}
	return
}

func MinMaxPtX(pts []image.Point) (minX, maxX float64) {
	minX, maxX = math.MaxInt32, math.MinInt32
	xs, _ := UnzipPoints(pts)
	for _, x := range xs {
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
	}
	return
}

func SortPointByX(pts []image.Point) {
	sort.SliceStable(pts, func(i, j int) bool {
		return (pts)[i].X < (pts)[j].X
	})
}

func SortPointByY(pts []image.Point) {
	sort.SliceStable(pts, func(i, j int) bool {
		return (pts)[i].Y < (pts)[j].Y
	})
}

// func DrawLines(img image.Image, pts []image.Point, step int, savePath string) image.Image {

// 	dc := gg.NewContextForImage(img)
// 	dc.SetRGB(1, 0, 0)
// 	dc.SetLineWidth(2)

// 	SortPointByY(pts)
// 	for i := 1; i < len(pts); i++ {

// 	}
// }

func DrawSpline(img image.Image, pts []image.Point, step int, color string, savePath string) image.Image {

	dc := gg.NewContextForImage(img)

	switch color {
	case "RED", "R", "Red", "red":
		dc.SetRGB(1, 0, 0)
	case "GREEN", "G", "Green", "green":
		dc.SetRGB(0, 1, 0)
	case "BLUE", "B", "Blue", "blue":
		dc.SetRGB(0, 0, 1)
	default:
		dc.SetRGB(0, 0, 0)
	}

	dc.SetLineWidth(4)

	SortPointByX(pts)
	minX, maxX := MinMaxPtX(pts)
	maxX2 := maxX - float64(step)
	s := gospline.NewCubicSpline(UnzipPoints(pts))
	for x := minX; x <= maxX2; x += float64(step) {
		y := s.At(x)
		xNext := x + float64(step)
		yNext := s.At(xNext)
		dc.DrawLine(x, y, xNext, yNext)
	}
	dc.Stroke()

	if savePath != "" {
		dc.SavePNG(savePath)
	}

	return dc.Image()
}

func DrawCircle(img image.Image, centres []image.Point, r int, color string, savePath string) image.Image {
	dc := gg.NewContextForImage(img)
	for _, c := range centres {
		dc.DrawCircle(float64(c.X), float64(c.Y), float64(r))
	}
	switch color {
	case "RED", "R", "Red", "red":
		dc.SetRGB(1, 0, 0)
	case "GREEN", "G", "Green", "green":
		dc.SetRGB(0, 1, 0)
	case "BLUE", "B", "Blue", "blue":
		dc.SetRGB(0, 0, 1)
	default:
		dc.SetRGB(0, 0, 0)
	}
	dc.Fill()

	if savePath != "" {
		dc.SavePNG(savePath)
	}

	return dc.Image()
}
