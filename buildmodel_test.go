package main

import (
	"fmt"
	"image"
	"testing"
)

func TestBuildModel(t *testing.T) {
	BuildModelROI()
}

func TestFindColorArea(t *testing.T) {
	MarkedROI2JSON("./cfg/model-roi.jpg", "./cfg/key-area.json")
}

func TestLoadKeyArea(t *testing.T) {
	// load & remark
	roi := loadImage("./cfg/roi-(2500,50)-(3000,2162).jpg")
	mPt := LoadROI("./cfg/key-area.json")
	pts := []image.Point{}
	for pt := range mPt {
		fmt.Println(pt)
		pts = append(pts, pt)
	}
	DrawCircle(roi, pts, 1, "R", "./cfg/re-mark.jpg")
}
