package main

import "image"

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

func BuildModel() {
	roi := NewROI("./cfg/model.jpg", 2500, 50, 3000, 2162)
	saveJPG(roi.data, "./cfg/model-roi.jpg")
}
