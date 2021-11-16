package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
)

type ROI struct {
	left   int
	top    int
	right  int
	bottom int
	data   image.Image
}

func NewROI(file string, left, top, right, bottom int) ROI {
	std := loadImage(file) // 3842*2162
	return ROI{
		left:   left,
		top:    top,
		right:  right,
		bottom: bottom,
		data:   ROIrgba(std, left, top, right, bottom),
	}
}

func NewROI2(file string, rect image.Rectangle) ROI {
	return NewROI(file, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
}

func BuildModelROI() {
	rect := image.Rectangle{
		Min: image.Point{
			X: 2500,
			Y: 50,
		},
		Max: image.Point{
			X: 3000,
			Y: 2162,
		},
	}
	roi := NewROI2("./cfg/model.jpg", rect)
	saveJPG(roi.data, fmt.Sprintf("./cfg/roi-%v.jpg", rect))
}

func MarkedROI2JSON(mROIFile, outJSONFile string) {
	// search
	img := loadImage(mROIFile)
	mPt := FindColorArea(img, color.RGBA{0, 255, 0, 0})
	fmt.Println(len(mPt))

	// store
	mPtStore := map[string]interface{}{}
	for pt := range mPt {
		mPtStore[pt.String()] = "G"
	}
	data, err := json.Marshal(mPtStore)
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(outJSONFile, data, os.ModePerm)
}

func LoadROI(roijson string) map[image.Point]struct{} {

	// load
	data, err := os.ReadFile(roijson)
	if err != nil {
		log.Fatalln(err)
	}
	mPtLoad := map[string]interface{}{}
	err = json.Unmarshal(data, &mPtLoad)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(len(mPtLoad))

	//////////////////////////////////

	mPt := make(map[image.Point]struct{})
	for strPt := range mPtLoad {
		pt := image.Point{}
		fmt.Sscanf(strPt, "(%d,%d)", &pt.X, &pt.Y)
		mPt[pt] = struct{}{}
	}

	return mPt
}
