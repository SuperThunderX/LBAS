package main

import (
	"image"
	"image/draw"
)

func ROIrgba(img image.Image, left, top, right, bottom int) *image.RGBA {
	rect := image.Rect(0, 0, right-left, bottom-top)
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, image.Point{left, top}, draw.Src)
	return rgba
}

func ROIrgbaV2(img image.Image, cx, cy, sRadius int) *image.RGBA {
	left := cx - sRadius
	top := cy - sRadius
	right := cx + sRadius
	bottom := cy + sRadius
	return ROIrgba(img, left, top, right, bottom)
}

func ROIcmyk(img image.Image, left, top, right, bottom int) *image.CMYK {
	rect := image.Rect(0, 0, right-left, bottom-top)
	cmyb := image.NewCMYK(rect)
	draw.Draw(cmyb, rect, img, image.Point{left, top}, draw.Src)
	return cmyb
}

func ROIcmykV2(img image.Image, cx, cy, sRadius int) *image.CMYK {
	left := cx - sRadius
	top := cy - sRadius
	right := cx + sRadius
	bottom := cy + sRadius
	return ROIcmyk(img, left, top, right, bottom)
}

func ROIgray(img image.Image, left, top, right, bottom int) *image.Gray {
	rect := image.Rect(0, 0, right-left, bottom-top)
	gray := image.NewGray(rect)
	draw.Draw(gray, rect, img, image.Point{left, top}, draw.Src)
	return gray
}

func ROIgrayV2(img image.Image, cx, cy, sRadius int) *image.Gray {
	left := cx - sRadius
	top := cy - sRadius
	right := cx + sRadius
	bottom := cy + sRadius
	return ROIgray(img, left, top, right, bottom)
}

func GrayAve(img *image.Gray) byte {
	sum, n := 0, 0
	rect := img.Bounds()
	for y := 0; y < rect.Dy(); y++ {
		pHead := img.Pix[y*img.Stride:]
		for x := 0; x < rect.Dx(); x++ {
			pxl := pHead[x]
			sum += int(pxl)
			n++
		}
	}
	return byte(sum / n)
}

func Cvt2Gray(img image.Image) *image.Gray {
	rect := img.Bounds()
	gray := image.NewGray(rect)
	draw.Draw(gray, rect, img, rect.Min, draw.Src)
	return gray
}

func GrayStripeV(img *image.Gray, x int) (stripe []int) {
	for y := 0; y < img.Rect.Dy(); y++ {
		offset := y * img.Stride
		pixel := img.Pix[offset+x]
		stripe = append(stripe, int(pixel))
	}
	return
}

func GrayStripeH(img *image.Gray, y int) (stripe []int) {
	offset := img.Stride * y
	line := img.Pix[offset : offset+img.Stride]
	for _, p := range line {
		stripe = append(stripe, int(p))
	}
	return
}
