package main

import (
	"fmt"
	"image"
	"math"
	"sort"
)

func absDiffInt(a, b int) int {
	return int(math.Abs(float64(a) - float64(b)))
}

func absDiffByte(a, b byte) byte {
	return byte(math.Abs(float64(a) - float64(b)))
}

type ROICandidate struct {
	pt   image.Point
	data *image.RGBA
	dy   int
}

func nextROICandidates(imgPath, cfgEdge string) (selected []ROICandidate) {
	img := loadImg(imgPath)
	edge := LoadLastRecord(cfgEdge)
	for _, pt := range edge.Points() {
		selected = append(selected, ROICandidate{
			pt:   pt,
			data: ROIrgbaV2(img, pt.X, pt.Y, roiRadius),
		})
	}
	sort.SliceStable(selected, func(i, j int) bool {
		return selected[i].pt.X < selected[j].pt.X
	})
	return
}

func makeNextEdgeCfg(selected []ROICandidate, cfgEdge, recordName, imgPath string, slopeStep int) {
	record := NewEdgeRecord(recordName, imgPath)
	for _, roi := range selected {
		f := feature(roi.data, slopeStep)
		above, below, left, right := f[0], f[1], f[2], f[3]
		record.AddPtInfo(roi.pt.X, roi.pt.Y, above, below, left, right)
	}
	record.Log(cfgEdge)
}

func NextKeyPoints(imgPath, cfgEdge, nextRecordName string, slopeStep int) (centres []image.Point) {

	img := loadImg(imgPath)
	edge := LoadLastRecord(cfgEdge)
	selected := nextROICandidates(imgPath, cfgEdge)

	pts4all := []ROICandidate{}

	for _, roi := range selected {

		pts4each := []ROICandidate{}

		// looking for edge config
		for _, pt := range edge.Pts {

			// refer to suitable config roi
			if roi.pt.X == pt.X {

				// gray := Cvt2Gray(roi.data)
				r, _, _, _ := SplitRGBA(roi.data) // choose [red] channel for slope
				ptr := GrayStripeV(r, roiRadius)
				ps := slope(ptr, slopeStep, 0)

				// if pt.ValAbove < pt.ValBelow {
				// 	// up -> down : dark -> bright
				// 	for _, s := range ps[:5] {
				// 	}
				// }

				if pt.ValAbove > pt.ValBelow {
					// up -> down : bright -> dark
					for i := len(ps) - 1; i >= len(ps)-5; i-- {
						s := ps[i]

						y := roi.pt.Y - roiRadius + s
						tempROI := ROIrgbaV2(img, pt.X, y, roiRadius)
						f := feature(tempROI, slopeStep)
						above, below := f[0], f[1]

						if absDiffByte(pt.ValAbove, above) < ERR && absDiffByte(pt.ValBelow, below) < ERR {
							pts4each = append(pts4each, ROICandidate{
								pt: image.Point{
									X: pt.X,
									Y: y,
								},
								data: tempROI,
								dy:   absDiffInt(y, pt.Y),
							})
						}
					}
				}
			}
		}

		sort.SliceStable(pts4each, func(i, j int) bool {
			return pts4each[i].dy < pts4each[j].dy
		})

		if len(pts4each) > 0 {
			wanted := pts4each[0]
			centres = append(centres, wanted.pt)
			pts4all = append(pts4all, ROICandidate{
				pt:   wanted.pt,
				data: wanted.data,
			})
		}
	}

	////////////////////////////////////////////

	nPoor := 0
	for _, ctr := range centres {
		for _, pt := range edge.Pts {
			// refer to suitable config roi
			if ctr.X == pt.X {
				if absDiffInt(ctr.Y, pt.Y) > roiRadius/2 {
					nPoor++
				}
			}
		}
	}
	if nPoor >= len(selected)/2 {
		fmt.Println("border curve may not be accurate @", imgPath)	
		return
	}

	////////////////////////////////////////////

	// [pts4all] for next config
	makeNextEdgeCfg(pts4all, cfgEdge, nextRecordName, imgPath, slopeStep)
	return
}
