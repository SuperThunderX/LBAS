package main

func main() {

	// inFolder := flag.String("in", "./in", "images folder path for batch input")
	// stdImage := flag.String("std", "./in/std.jpg", "path of marked starting image")
	// outFolder := flag.String("out", "./out", "path for output")
	// flag.Parse()

	// switch {
	// case !filedir.FileExists(*stdImage):
	// 	fmt.Println("use '-std path/to/std.jpg' to set marked starting image")
	// 	return
	// case !filedir.DirExists(*inFolder):
	// 	fmt.Println("use '-in path/to/images/folder' to set images folder")
	// 	return
	// }

	// ///////////////////////////////////////////////////////////

	// os.Mkdir("./cfg/", os.ModePerm)

	// cfgEdgeAB := "./cfg/AB-edge.json"
	// cfgEdgeBC := "./cfg/BC-edge.json"

	// BuildModel(cfgEdgeAB, "AB", *stdImage, color.RGBA{255, 0, 0, 255}, 7)
	// BuildModel(cfgEdgeBC, "BC", *stdImage, color.RGBA{0, 255, 0, 255}, 11)

	// ///////////////////////////////////////////////////////////

	// mode := "LINE"

	// var draw func(img image.Image, centres []image.Point, r int, color string, savePath string) image.Image

	// switch mode {
	// case "DOT":
	// 	draw = DrawCircle
	// case "LINE":
	// 	draw = DrawSpline
	// default:
	// 	draw = DrawSpline
	// }

	// /////////////////////////////////////////////

	// files, _, err := filedir.WalkFileDir(*inFolder, false)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// sort.SliceStable(files, func(i, j int) bool {
	// 	return files[i] < files[j]
	// })

	// for _, file := range files {

	// 	if file == *stdImage {
	// 		continue
	// 	}

	// 	fmt.Println(file)

	// 	inImage := file
	// 	outImage := filepath.Join(*outFolder, filepath.Base(file))

	// 	cfgEdge := "./cfg/AB-edge.json"
	// 	color := "R"
	// 	rs := 5 // *** [dot-radius] or [line-step] ***

	// 	img := loadImg(inImage)
	// 	keyPts := NextKeyPoints(inImage, cfgEdge, "", 7)
	// 	draw(img, keyPts, rs, color, outImage)

	// 	// cfgEdge = "./cfg/BC-edge.json"
	// 	// color = "G"

	// 	// inImage = outImage
	// 	// img = loadImg(inImage)
	// 	// keyPts = NextKeyPoints(inImage, cfgEdge, "", 11)
	// 	// draw(img, keyPts, rs, color, inImage)

	// 	fmt.Println()
	// }

}
