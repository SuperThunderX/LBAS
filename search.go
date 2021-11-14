package main

import (
	"fmt"
	"image"
	"image/color"

	gocv "github.com/digisan/go-handy-cv/blob"
	gotkio "github.com/digisan/gotk/io"
)

func FindPosByClr(img image.Image, c color.RGBA) (pos []image.Point) {
	rect := img.Bounds()
	rgba := ROIrgba(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
	for y := 0; y < rect.Dy(); y++ {
		start := y * rgba.Stride
		pln := rgba.Pix[start : start+rgba.Stride]
		for x := 0; x < rect.Dx(); x++ {
			p := pln[4*x:]
			cmp := color.RGBA{p[0], p[1], p[2], p[3]}
			if ColorEqual(c, cmp, 3, 3, 3, 256) {
				pos = append(pos, image.Point{x, y})
			}
		}
	}
	return
}

func FindROIrgbaByClr(img image.Image, c color.RGBA, sRadius, iRadius int, auditPath string) (mPtROI map[image.Point]*image.RGBA) {

	mPtRGBA := make(map[image.Point]*image.RGBA)
	for _, pos := range FindPosByClr(img, c) {
		roi := ROIrgbaV2(img, pos.X, pos.Y, sRadius)
		mPtRGBA[pos] = roi
	}

	mPtROI = make(map[image.Point]*image.RGBA)
NEXT:
	for pt1, rgba := range mPtRGBA {
		for pt2 := range mPtROI {
			if gocv.PtDis(pt1, pt2) < iRadius {
				continue NEXT
			}
		}
		mPtROI[pt1] = rgba
	}

	I := 0
	for pt, roi := range mPtROI {
		if len(auditPath) > 0 {
			gotkio.MustCreateDir(auditPath)
			savePNG(roi, fmt.Sprintf("./%s/%00d-%d-%d.png", "./out/audit/", I, pt.X, pt.Y))
			I++
		}
	}

	return
}

// func FindROIrgbaByBlob(img image.Image,
// 	sRadius int,
// 	filterR func(x, y int, p byte) bool,
// 	filterG func(x, y int, p byte) bool,
// 	filterB func(x, y int, p byte) bool,
// 	disErr int,
// 	auditPath string) (mPtRGBA map[image.Point]*image.RGBA) {

// 	gotkio.MustCreateDir(auditPath)
// 	mPtRGBA = make(map[image.Point]*image.RGBA)

// 	// rect := img.Bounds()
// 	// rgba := ROIrgba(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)

// 	r, g, b, _ := SplitRGBA(img)
// 	blobPosGrp := gocv.DetectClrBlobPos(r.Rect.Dx(), r.Rect.Dy(), r.Stride,
// 		r.Pix, g.Pix, b.Pix,
// 		filterR, filterG, filterB, disErr)

// 	for i, bpos := range blobPosGrp {
// 		pos := image.Point{X: bpos.X, Y: bpos.Y}
// 		roi := ROIrgbaV2(img, pos.X, pos.Y, sRadius)
// 		mPtRGBA[pos] = roi
// 		if len(auditPath) > 0 {
// 			savePNG(roi, fmt.Sprintf("./%s/%00d-%d-%d.png", auditPath, i, pos.X, pos.Y))
// 		}
// 	}
// 	return
// }
