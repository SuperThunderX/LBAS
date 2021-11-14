package main

import (
	"fmt"
	"image"
	"image/color"
)

const (
	searchOffset  = 40
	roiRadius     = 40
	ignrClrRadius = 5
	ERR           = 15
)

const (
	iGray = iota
	iRed
	iGreen
	iBlue
	iN
)

var (
	chClr  = []string{"Gray", "Red", "Green", "Blue"}
	wChAve = []float64{0.1, 0.8, 0.1, 0.0}
)

func weight(values [4]byte) byte {
	sum := 0.0
	for i := 0; i < 4; i++ {
		sum += float64(values[i]) * wChAve[i]
	}
	return byte(sum)
}

func feature(roi image.Image, slopeStep int) [4]byte {

	gray := Cvt2Gray(roi)
	r, g, b, _ := SplitRGBA(roi)

	aboves, belows, lefts, rights := [4]byte{}, [4]byte{}, [4]byte{}, [4]byte{}

	wSubRadius := slopeStep / 3

	for iCh, ch := range []*image.Gray{gray, r, g, b} {

		aboveSubPt := roiRadius - slopeStep/2
		sub := ROIgrayV2(ch, roiRadius, aboveSubPt, wSubRadius)
		aboves[iCh] = GrayAve(sub)

		belowSubPt := roiRadius + slopeStep/2
		sub = ROIgrayV2(ch, roiRadius, belowSubPt, wSubRadius)
		belows[iCh] = GrayAve(sub)

		leftSubPt := roiRadius - slopeStep/2
		sub = ROIgrayV2(ch, leftSubPt, roiRadius, wSubRadius)
		lefts[iCh] = GrayAve(sub)

		rightSubPt := roiRadius + slopeStep/2
		sub = ROIgrayV2(ch, rightSubPt, roiRadius, wSubRadius)
		rights[iCh] = GrayAve(sub)
	}

	return [4]byte{weight(aboves), weight(belows), weight(lefts), weight(rights)}
}

func BuildModel(recordPath, recordName, imagePath string, aim color.RGBA, slopeStep int) {

	record := NewEdgeRecord(recordName, imagePath)
	img := loadImg(imagePath)
	mPtROI := FindROIrgbaByClr(img, aim, roiRadius, ignrClrRadius, "./out/")
	// fmt.Println(mPtROI)

	for pt, roi := range mPtROI {
		f := feature(roi, slopeStep)
		above, below, left, right := f[0], f[1], f[2], f[3]
		record.AddPtInfo(pt.X, pt.Y, above, below, left, right)
		fmt.Println()
	}

	record.Log(recordPath)
}
