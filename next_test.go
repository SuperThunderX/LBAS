package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/digisan/gotk/filedir"
)

func TestSearchNextROI(t *testing.T) {

	os.Mkdir("./cfg/", os.ModePerm)

	modelImage := "./in/11.02/std.jpg"

	cfgEdgeAB := "./cfg/AB-edge.json"
	cfgEdgeBC := "./cfg/BC-edge.json"

	BuildModel(cfgEdgeAB, "AB", modelImage, color.RGBA{255, 0, 0, 255}, 7)
	BuildModel(cfgEdgeBC, "BC", modelImage, color.RGBA{0, 255, 0, 255}, 11)

	///////////////////////////////////////////////////////////

	mode := "LINE"

	var draw func(img image.Image, centres []image.Point, r int, color string, savePath string) image.Image

	switch mode {
	case "DOT":
		draw = DrawCircle
	case "LINE":
		draw = DrawSpline
	default:
		draw = DrawSpline
	}

	/////////////////////////////////////////////

	files, _, err := filedir.WalkFileDir("./in/11.02", false)
	if err != nil {
		log.Fatalln(err)
	}

	sort.SliceStable(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	for _, file := range files {

		if strings.HasSuffix(file, "std.jpg") {
			continue
		}

		fmt.Println(file)

		inImage := file
		outImage := filepath.Join("./out/", filepath.Base(file))

		cfgEdge := "./cfg/AB-edge.json"
		color := "R"
		rs := 5 // *** [dot-radius] or [line-step] ***

		img := loadImg(inImage)
		keyPts := NextKeyPoints(inImage, cfgEdge, "", 7)
		draw(img, keyPts, rs, color, outImage)

		cfgEdge = "./cfg/BC-edge.json"
		color = "G"

		inImage = outImage
		img = loadImg(inImage)
		keyPts = NextKeyPoints(inImage, cfgEdge, "", 11)
		draw(img, keyPts, rs, color, inImage)

		fmt.Println()
	}
}
