package main

import (
	"flag"
	"fmt"
	"image"
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

	caModePtr := flag.Bool("ca", false, "creating detect area mode")
	caImgPtr := flag.String("caimg", "./in/sample.jpg", "input marked image")
	caClrPtr := flag.String("caclr", "R", "marked color, (R G B W K C Y M)")
	caOutPtr := flag.String("caout", "./cfg/sample.json", "output json file")

	inpathPtr := flag.String("in", "./in/sample", "input folder for original images")
	outpathPtr := flag.String("out", "./out", "output folder for scaled images")
	cfgpathPtr := flag.String("cfg", "./cfg/config.json", "config file(json) for running program")
	
	flag.Parse()

	if *caModePtr {
		MarkedArea2JSON(*caImgPtr, *caOutPtr, *caClrPtr)
		return
	}

	defer gotk.TrackTime(time.Now())

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

		marks := [][]image.Point{}
		output := filepath.Join(outpath, filename)
		results := []string{}

		for _, cfgArea := range cfg.Roi.AimAreas {

			// audit marked ROI image
			// markedFile := filepath.Join(outpath, fmt.Sprintf("%s-roi-mark.jpg", filename))
			// MarkAreaFromJSON(roiFile, cfgArea.File, cfgArea.Offset, 0, markedFile, "R")

			N := 100
			areas := SplitAreaPts(cfgArea.File, cfgArea.Offset, 0, "Y", N)
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

			// detectFrom := int(float64(N) * cfg.ValidRange[0])
			detectTo := int(float64(N) * cfgArea.ValidRange[1])
			thBright := cfgArea.ThBright
			thContrast := cfgArea.ThContrast

			for i := 3; i < detectTo; i++ {
				p := i
				ave := aves[i]
				if aves[p-1] > thBright || aves[p-2] > thBright || aves[p-3] > thBright {
					fmt.Println("manual detect 1")
					if ave > aves[p+1] && ave > aves[p+2] && ave > aves[p+3] {
						fmt.Println("manual detect 2")
						if absDiffFloat(ave, aves[p+1]) > thContrast &&
							absDiffFloat(ave, aves[p+2]) > thContrast &&
							absDiffFloat(ave, aves[p+3]) > thContrast {
							fmt.Println("manual detect 3")
							if absDiffFloat(aves[p+1], aves[p+2]) < thContrast*2 &&
								absDiffFloat(aves[p+2], aves[p+3]) < thContrast*2 {
								fmt.Println("manual detect 4")
								Indices = append(Indices, p)
								break
							}
						}
					}
				}
			}

			fmt.Printf("check 1 done [%d]\n", len(Indices))

			if len(Indices) == 0 {
				fmt.Println("auto detect")

				for _, s3 := range scale3list {
					p := s3.index
					ave := aves[p]
					if p >= 3 && p <= len(scale3list) {
						// *** up bright, down dark ***
						// normal image
						if ave > aves[p+1] && ave > aves[p+2] && ave > aves[p+3] {
							if aves[p-3] > ave && aves[p-2] > ave && aves[p-1] > ave {
								if p > 10 && p < len(aves)-10 {
									if absDiffFloat(ave, aves[p+9]) < 30 &&
										absDiffFloat(aves[p-9], aves[0]) < 20 &&
										absDiffFloat(aves[p+9], aves[0]) > 50 {
										Indices = append(Indices, p)
										break
									}
								}
							}
						}
					}
				}
			}

			if len(Indices) == 0 {
				lk.Log2F(true, "./ignore.log")
				lk.Log("Ignore 2, Cannot find scale, @%s", file)
				continue // next area
			}

			Index := Indices[0]
			scale := float64(Index) / float64(len(areas))

			///

			a := areas[Index][0]
			fmt.Println("output:", a)

			///

			ymin, ymax := cfgArea.ValidRange[0], cfgArea.ValidRange[1]
			if scale >= ymin && scale < ymax {
				// audit scaled ROI image
				marks = append(marks, areas[Index])
				results = append(results, fmt.Sprintf("%.02f", scale))
				// MarkArea(roiFile, output, "R", areas[Index])
			}

		} // end of each aim-area

		// draw all
		all := []image.Point{}
		for _, mark := range marks {
			all = append(all, mark...)
		}
		MarkArea(roiFile, output+".jpg", "R", all)

		// write results
		r := strings.Join(results, "\n")
		os.WriteFile(filepath.Join(filepath.Dir(output), filename+".txt"), []byte(r), os.ModePerm)

	} // end of each file
}
