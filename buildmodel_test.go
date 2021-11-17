package main

import (
	"fmt"
	"image"
	"testing"
)

func TestModelROI(t *testing.T) {
	imgfile := "./in/sample.jpg"
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
	roifile := fmt.Sprintf("./cfg/roi-%v.jpg", rect)
	ModelROI(imgfile, roifile, rect)
}

func TestFindColorArea(t *testing.T) {
	// Generate Area Configuration
	MarkedArea2JSON("./cfg/model-roi.jpg", "./cfg/key-area.json")
}

func TestLoadKeyArea(t *testing.T) {
	// Audit marked ROI image
	MarkAreaFromJSON("./cfg/model-roi.jpg", "./cfg/key-area.json", "./cfg/re-mark.jpg", "B")
}

func TestSplitAreaPts(t *testing.T) {
	areas := SplitAreaPts("./cfg/key-area.json", "Y", 100)
	for i, area := range areas {
		fmt.Println(i, len(area))
	}
}

func TestStatAreaPixel(t *testing.T) {
	areas := SplitAreaPts("./cfg/key-area.json", "Y", 100)
	img := loadImage("./cfg/roi-(2500,50)-(3000,2162).jpg")

	aves := []float64{}
	for _, area := range areas {
		a, _ := StatAreaPixel(img, area, 0.65, 0.35, 0)
		// fmt.Printf("%03d -- %3.02f -- %2.02f\n", i, a, s)
		aves = append(aves, a)
	}

	for i := 0; i < len(aves)-3; i++ {
		batch := aves[i : i+3]
		a := AveFloat(batch)
		s := StdDevFloat(batch)
		fmt.Printf("%d - %.2f - %.2f\n", i, a, s)
	}

	MarkArea("./cfg/roi-(2500,50)-(3000,2162).jpg", "./out/marked.jpg", "G", areas[1])
}
