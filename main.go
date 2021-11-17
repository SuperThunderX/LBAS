package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/digisan/gotk"
)

func main() {

	defer gotk.TrackTime(time.Now())

	if err := os.MkdirAll("./temp", os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	cfg := LoadConfig("./cfg/config.json")
	fmt.Println(cfg)

	roiFile := "./temp/sample-roi.jpg"
	ModelROI("./in/sample.jpg", roiFile, cfg.ROIRect())

	// audit marked ROI image
	// markedFile := "./temp/sample-roi-mark.jpg"
	// MarkAreaFromJSON(roiFile, cfg.KeyArea, markedFile, "R")

	N := 100
	areas := SplitAreaPts(cfg.KeyArea, "Y", N)
	img := loadImage(roiFile)

	aves := []float64{}
	for _, area := range areas {
		a, _ := StatAreaPixel(img, area, 0.6, 0.4, 0)
		// fmt.Printf("%d - %.2f - %.2f\n", i, a, sd)
		aves = append(aves, a)
	}

	fmt.Println("-------------------------------------")

	type scale3 struct {
		index int
		ave   float64
		sd    float64
	}

	scale3list := []scale3{}

	for i := 0; i < len(aves)-3; i++ {
		batch := aves[i : i+3]
		ave := AveFloat(batch)
		sd := StdDevFloat(batch)
		// fmt.Printf("%d - %.2f - %.2f\n", i, ave, sd)

		scale3list = append(scale3list, scale3{
			index: i,
			ave:   ave,
			sd:    sd,
		})
	}

	sort.SliceStable(scale3list, func(i, j int) bool {
		return scale3list[i].sd > scale3list[j].sd
	})

	Indices := []int{}
	for _, s3 := range scale3list {
		p := s3.index
		ave := aves[p]
		if p >= 3 && p <= len(scale3list) {
			// *** up bright, down dark ***
			if ave > aves[p+1] && ave > aves[p+2] && ave > aves[p+3] {
				if aves[p-3] > ave && aves[p-2] > ave && aves[p-1] > ave {
					Indices = append(Indices, p)
				}
			}
		}
	}

	// sort.SliceStable(Indices, func(i, j int) bool {
	// 	return Indices[i] < Indices[j]
	// })

	Index := Indices[0]
	scale := float64(Index) / float64(len(areas))

	// audit scaled ROI image
	scaleFile := fmt.Sprintf("./temp/sample-roi-%.02f.jpg", scale)
	MarkArea(roiFile, scaleFile, "G", areas[Index])

}
