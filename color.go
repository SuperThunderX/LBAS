package main

import (
	"image"
	"image/color"
	"log"
)

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

// func SplitHSV(img image.Image) (h, s, v *image.Gray) {
// 	// rect := img.Bounds()
// 	// rgba := ROIrgba(img, rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
// 	// h, s, v = image.NewGray(rect), image.NewGray(rect), image.NewGray(rect)
// }
