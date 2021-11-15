package main

import (
	// _ "image/jpg"
	"encoding/json"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"testing"
	"time"

	"github.com/digisan/gotk"
)

func TestSplitRGBA(t *testing.T) {
	defer gotk.TrackTime(time.Now())

	img := loadImage("./in/sample/1.jpg")
	fmt.Println(img.Bounds())

	r, g, b, a := SplitRGBA(img)
	saveJPG(r, "./out/r.jpg")
	saveJPG(g, "./out/g.jpg")
	saveJPG(b, "./out/b.jpg")
	saveJPG(a, "./out/a.jpg")

	// ///

	com := CompositeRGBA(r, g, b, a)
	saveJPG(com, "./out/com1.png")
}

func TestFindColorArea(t *testing.T) {

	// search
	img := loadImage("./cfg/model-roi.jpg")
	mPt := FindColorArea(img, color.RGBA{0, 255, 0, 0})
	fmt.Println(len(mPt))

	// store
	mPtStore := map[string]interface{}{}
	for pt := range mPt {
		mPtStore[pt.String()] = "G"
	}
	data, err := json.Marshal(mPtStore)
	if err != nil {
		log.Fatalln(err)
	}
	os.WriteFile("./cfg/key-area.json", data, os.ModePerm)

	// load
	data, err = os.ReadFile("./cfg/key-area.json")
	if err != nil {
		log.Fatalln(err)
	}
	mPtLoad := map[string]interface{}{}
	err = json.Unmarshal(data, &mPtLoad)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(len(mPtLoad))
}
