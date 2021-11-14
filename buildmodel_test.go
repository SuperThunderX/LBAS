package main

import (
	"fmt"
	"image/color"
	"os"
	"testing"
)

func TestBuildModel(t *testing.T) {
	BuildModel("./out/AB-edge.json", "AB", "./in/start/1.jpg", color.RGBA{255, 0, 0, 255}, 7)
}

func TestLoadEdgeRecord(t *testing.T) {
	edge := LoadLastRecord("./out/AB-edge.json")
	fmt.Println(edge)
	fmt.Println(edge.Pts[0].ValAbove)
}

func TestDrawEdge(t *testing.T) {

	os.Mkdir("./cfg/", os.ModePerm)

	modelImage := "./in/sample/std.jpg"

	cfgEdgeAB := "./cfg/AB-edge.json"
	cfgEdgeBC := "./cfg/BC-edge.json"

	BuildModel(cfgEdgeAB, "AB", modelImage, color.RGBA{255, 0, 0, 255}, 7)
	BuildModel(cfgEdgeBC, "BC", modelImage, color.RGBA{0, 255, 0, 255}, 11)

	// edge := LoadLastRecord(cfgEdge)

	// outImage := "./out/lines.jpg"
	// img := loadImg(modelImage)
	// DrawSpline(img, edge.Points(), 5, outImage)
}

func TestMarkImage(t *testing.T) {

	cfgEdge := "./cfg/AB-edge.json"
	edge := LoadLastRecord(cfgEdge)

	inImage := "./in/sample/1.jpg"
	outImage := "./out/1.jpg"
	img := loadImg(inImage)
	DrawSpline(img, edge.Points(), 5, "B", outImage)

	// inImage = "./in/sample/2.jpg"
	// outImage = "./out/2.jpg"
	// img = loadImg(inImage)
	// DrawSpline(img, edge.Points(), 5, outImage)

	// inImage = "./in/sample/3.jpg"
	// outImage = "./out/3.jpg"
	// img = loadImg(inImage)
	// DrawSpline(img, edge.Points(), 5, outImage)

	// inImage = "./in/sample/4.jpg"
	// outImage = "./out/4.jpg"
	// img = loadImg(inImage)
	// DrawSpline(img, edge.Points(), 5, outImage)
}
