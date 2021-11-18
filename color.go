package main

import (
	"image"
	"image/color"
	"log"
	"math"
)

const (
	ClrErr = 4
)

func Color2RGBA(clr color.Color) (r, g, b, a byte) {
	r32, g32, b32, a32 := clr.RGBA()
	return byte(r32 >> 8), byte(g32 >> 8), byte(b32 >> 8), byte(a32 >> 8)
}

func Color2Gray(clr color.Color) byte {
	gray32, _, _, _ := clr.RGBA()
	return byte(gray32 >> 8)
}

func absDiffInt(a, b int) int {
	return int(math.Abs(float64(a) - float64(b)))
}

func absDiffByte(a, b byte) byte {
	return byte(math.Abs(float64(a) - float64(b)))
}

func absDiffFloat(a, b float64) float64 {
	return math.Abs(a - b)
}

func ColorEqual(c1, c2 color.RGBA, eR, eG, eB, eA int) bool {
	if abs(int(c1.R)-int(c2.R)) <= eR &&
		abs(int(c1.G)-int(c2.G)) <= eG &&
		abs(int(c1.B)-int(c2.B)) <= eB &&
		abs(int(c1.A)-int(c2.A)) <= eA {
		return true
	}
	return false
}

func CompositeRGBA(r, g, b, a image.Image) *image.RGBA {

	rectR, rectG, rectB, rectA := r.Bounds(), g.Bounds(), b.Bounds(), a.Bounds()
	if rectR != rectG || rectG != rectB || rectB != rectA {
		log.Fatalln("r, g, b, a all must be same size")
		return nil
	}

	rgba := image.NewRGBA(rectR)
	bytes := rgba.Pix
	for i, p := range r.(*image.Gray).Pix {
		bytes[i*4] = p
	}
	for i, p := range g.(*image.Gray).Pix {
		bytes[i*4+1] = p
	}
	for i, p := range b.(*image.Gray).Pix {
		bytes[i*4+2] = p
	}
	for i, p := range a.(*image.Gray).Pix {
		bytes[i*4+3] = p
	}
	return rgba
}

func FindColorArea(img image.Image, clr color.RGBA) map[image.Point]struct{} {

	// out
	mPt := make(map[image.Point]struct{})

	r, g, b, _ := SplitRGBA(img)
	rect := img.Bounds()

	mChClr := map[int]byte{
		0: clr.R,
		1: clr.G,
		2: clr.B,
	}
	mPtN := make(map[image.Point]int)
	for i, ch := range []*image.Gray{r, g, b} {
		for y := 0; y < rect.Dy(); y++ {
			offset := y * ch.Stride
			line := ch.Pix[offset:]
			for x := 0; x < rect.Dx(); x++ {
				pxl := line[x]
				pt := image.Point{X: x, Y: y}
				if absDiffByte(pxl, mChClr[i]) < ClrErr {
					mPtN[pt]++
				}
			}
		}
	}

	for pt, n := range mPtN {
		if n == 3 {
			mPt[pt] = struct{}{}
		}
	}

	return mPt
}

func SplitRGBA(img image.Image) (r, g, b, a *image.Gray) {

	rect := img.Bounds()

	left, top, right, bottom := rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y
	img = ROIrgba(img, left, top, right, bottom)

	var bytes []byte
	switch pImg := img.(type) {
	case *image.RGBA:
		bytes = pImg.Pix
	case *image.NRGBA:
		bytes = pImg.Pix
	// case *image.YCbCr: //	YCbCrSubsampleRatio444
	// 	bytes = pImg.Pix
	default:
		log.Fatalf("[%v] is not support", pImg)
	}

	r, g, b, a = image.NewGray(rect), image.NewGray(rect), image.NewGray(rect), image.NewGray(rect)
	for i, p := range bytes {
		switch i % 4 {
		case 0:
			r.Pix[i/4] = p
		case 1:
			g.Pix[i/4] = p
		case 2:
			b.Pix[i/4] = p
		case 3:
			a.Pix[i/4] = p
		}
	}

	// wg := &sync.WaitGroup{}
	// wg.Add(4)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[0:], r.Pix)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[1:], g.Pix)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[2:], b.Pix)
	// go func(rgbaBytes, chBytes []byte) {
	// 	for i, j := 0, 0; i < len(rgbaBytes); i += 4 {
	// 		chBytes[j] = rgbaBytes[i]
	// 		j++
	// 	}
	// 	wg.Done()
	// }(rgba.Pix[3:], a.Pix)
	// wg.Wait()

	return
}
