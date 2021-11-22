package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/digisan/gotk"
	"github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
)

func main() {

	defer gotk.TrackTime(time.Now())

	inpathPtr := flag.String("in", "./in/sample", "input folder for original images")
	outpathPtr := flag.String("out", "./out", "output folder for scaled images")
	cfgpathPtr := flag.String("cfg", "./cfg/config.json", "config file(json) for running program")
	flag.Parse()

	inpath := *inpathPtr
	outpath := *outpathPtr
	cfgpath := *cfgpathPtr

	if err := os.MkdirAll(outpath, os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	cfg := LoadConfig(cfgpath)
	fmt.Println(cfg)

	files, _, err := filedir.WalkFileDir(inpath, false)
	if err != nil {
		log.Fatalln(err)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	for _, file := range files {

		// if i > 10 {
		// 	break
		// }

		inFile := file
		filename := strings.TrimSuffix(filepath.Base(inFile), ".jpg")
		roiFile := filepath.Join(outpath, fmt.Sprintf("%s-roi.jpg", filename))
		ModelROI(inFile, roiFile, cfg.ROIRect())

		// audit marked ROI image
		// markedFile := filepath.Join(outpath, fmt.Sprintf("%s-roi-mark.jpg", filename))
		// MarkAreaFromJSON(roiFile, cfg.KeyArea, cfg.KeyAreaOffset, 0, markedFile, "R")

		N := 100
		areas := SplitAreaPts(cfg.KeyArea, cfg.KeyAreaOffset, 0, "Y", N)
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

		if len(scale3list) == 0 {
			lk.Log2F(true, "./ignore.log")
			lk.Log("Ignore 1, Cannot find scale, @%s", file)
			continue
		}

		Indices := []int{}
		for _, s3 := range scale3list {
			p := s3.index
			ave := aves[p]
			if p >= 3 && p <= len(scale3list) {
				// *** up bright, down dark ***
				// normal image
				if ave > aves[p+1] && ave > aves[p+2] && ave > aves[p+3] {
					if aves[p-3] > ave && aves[p-2] > ave && aves[p-1] > ave {
						if p > 10 && p < len(aves)-10 {
							if absDiffFloat(ave, aves[p+9]) < 20 &&
								absDiffFloat(aves[p-9], aves[0]) < 20 &&
								absDiffFloat(aves[p+9], aves[0]) > 40 {
								Indices = append(Indices, p)
								continue
							}
						}
					}
				}
			}
		}

		// sort.SliceStable(Indices, func(i, j int) bool {
		// 	return Indices[i] < Indices[j]
		// })

		// check strong reflection area
		if len(Indices) == 0 {
			fmt.Println("extra detect 0")

			// detectFrom := int(float64(N) * cfg.ValidRange[0])
			detectTo := int(float64(N) * cfg.ValidRange[1])
			thBright := cfg.ThBright     // 50.0
			thContrast := cfg.ThContrast // 5.0

			for i := 3; i < detectTo; i++ {
				p := i
				ave := aves[i]
				if aves[p-1] > thBright || aves[p-2] > thBright || aves[p-3] > thBright {
					fmt.Println("extra detect 1")
					if ave > aves[p+1] && ave > aves[p+2] && ave > aves[p+3] {
						fmt.Println("extra detect 2")
						if absDiffFloat(ave, aves[p+1]) > thContrast &&
							absDiffFloat(ave, aves[p+2]) > thContrast &&
							absDiffFloat(ave, aves[p+3]) > thContrast {
							fmt.Println("extra detect 3")
							if absDiffFloat(aves[p+1], aves[p+2]) < thContrast*2 &&
								absDiffFloat(aves[p+2], aves[p+3]) < thContrast*2 {
								fmt.Println("extra detect 4")
								Indices = append(Indices, p)
								break
							}
						}
					}
				}
			}
		}

		if len(Indices) == 0 {
			lk.Log2F(true, "./ignore.log")
			lk.Log("Ignore 2, Cannot find scale, @%s", file)
			continue // next image
		}

		Index := Indices[0]
		scale := float64(Index) / float64(len(areas))

		ymin, ymax := cfg.ValidRange[0], cfg.ValidRange[1]
		if scale >= ymin && scale < ymax {
			// audit scaled ROI image
			scaleFile := filepath.Join(outpath, fmt.Sprintf("%s-roi-%.02f.jpg", filename, scale))
			MarkArea(roiFile, scaleFile, "G", areas[Index])
		}
	}
}
