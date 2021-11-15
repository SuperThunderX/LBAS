package main

import (
	"fmt"
	"image/color"
	"testing"
	"time"

	"github.com/digisan/gotk"
)

func TestFindPosByColor(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	img := loadImage("./in/start/1.jpg")
	mPtROI := FindROIrgbaByClr(img, color.RGBA{254, 0, 0, 255}, 70, 8, "./out/")
	fmt.Println(mPtROI)
}

// func TestFindROIrgbaByBlob(t *testing.T) {
// 	defer gotk.TrackTime(time.Now())
//  img := loadImg("./in/start/1.jpg")
// 	mPtRGBA := FindROIrgbaByBlob(img, 70,
// 		func(x, y int, p byte) bool {
// 			return p > 253
// 		}, func(x, y int, p byte) bool {
// 			return p < 2
// 		}, func(x, y int, p byte) bool {
// 			return p < 2
// 		}, 5, "./out/")
// 	fmt.Println(mPtRGBA)
// }
