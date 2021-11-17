package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"sort"
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

func ModelROI(imgfile, roifile string, rect image.Rectangle) {
	roi := NewROI2(imgfile, rect)
	saveJPG(roi.data, roifile)
}

func MarkedArea2JSON(mROIFile, outJSONFile string) {
	// search
	img := loadImage(mROIFile)
	mPt := FindColorArea(img, color.RGBA{0, 255, 0, 0})
	fmt.Println(len(mPt))

	// store
	mPtStore := map[string]interface{}{}
	for pt := range mPt {
		mPtStore[pt.String()] = ""
	}
	data, err := json.Marshal(mPtStore)
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile(outJSONFile, data, os.ModePerm)
}

func LoadArea(areaJson string) map[image.Point]struct{} {

	// load
	data, err := os.ReadFile(areaJson)
	if err != nil {
		log.Fatalln(err)
	}
	mPtLoad := map[string]interface{}{}
	err = json.Unmarshal(data, &mPtLoad)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println(len(mPtLoad))

	//////////////////////////////////

	mPt := make(map[image.Point]struct{})
	for strPt := range mPtLoad {
		pt := image.Point{}
		fmt.Sscanf(strPt, "(%d,%d)", &pt.X, &pt.Y)
		mPt[pt] = struct{}{}
	}

	return mPt
}

func MarkAreaFromJSON(imgFile, areaJsonFile, markImgFile, color string) {
	// load & re-mark
	mPt := LoadArea(areaJsonFile)
	pts := []image.Point{}
	for pt := range mPt {
		pts = append(pts, pt)
	}
	DrawCircle(loadImage(imgFile), pts, 1, color, markImgFile)
}

func MarkArea(imgFile, markImgFile, color string, pts []image.Point) {
	DrawCircle(loadImage(imgFile), pts, 1, color, markImgFile)
}

func SplitAreaPts(areaJson, xy string, n int) [][]image.Point {

	areas := make([][]image.Point, int(float64(n)*1.1))

	pts := []image.Point{}
	for pt := range LoadArea(areaJson) {
		pts = append(pts, pt)
	}

	min, max := 0, 0
	switch xy {
	case "X", "x":
		sort.SliceStable(pts, func(i, j int) bool {
			return pts[i].X < pts[j].X
		})
		min, max = pts[0].X, pts[len(pts)-1].X
	case "Y", "y":
		sort.SliceStable(pts, func(i, j int) bool {
			return pts[i].Y < pts[j].Y
		})
		min, max = pts[0].Y, pts[len(pts)-1].Y
	default:
		log.Fatalln("[xy] can only set 'X' or 'Y'")
	}

	span := (max - min) / n

	switch xy {
	case "X", "x":
		for _, pt := range pts {
			i := (pt.X - min) / span
			areas[i] = append(areas[i], pt)
		}
	case "Y", "y":
		for _, pt := range pts {
			i := (pt.Y - min) / span
			areas[i] = append(areas[i], pt)
		}
	}

	sum, n := 0, 0
	for _, area := range areas {
		cnt := len(area)
		if cnt == 0 {
			break
		}
		sum += cnt
		n++
	}
	ave := sum / n

	end := len(areas)
	for i, area := range areas {
		if absDiffInt(len(area), ave) > ave/3 {
			end = i
			break
		}
	}

	return areas[:end]
}

func StatAreaPixel(img image.Image, area []image.Point, wR, wG, wB float64) (ave, sd float64) {

	var (
		rect = img.Bounds()
		rgba = ROIrgba(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
		// gray = Cvt2Gray(img)
	)

	var (
		rData = []byte{}
		gData = []byte{}
		bData = []byte{}
		// grayData = []byte{}
	)

	for _, pt := range area {
		clr := rgba.At(pt.X, pt.Y)
		// gval := gray.At(pt.X, pt.Y)
		// fmt.Println(i, clr, gval)

		r, g, b, _ := Color2RGBA(clr)
		// gray := Color2Gray(gval)

		rData = append(rData, r)
		gData = append(gData, g)
		bData = append(bData, b)
	}

	ave = AveByte(rData)*wR + AveByte(gData)*wG + AveByte(bData)*wB
	sd = StdDevByte(rData)*wR + StdDevByte(gData)*wG + StdDevByte(bData)*wB
	return
}
