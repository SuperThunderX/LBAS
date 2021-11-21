# LBAS

Liquid Boundary Analysis System

Compute & Mark 1st liquid border, computed scale value (0.0-1.0) is added on output image file name.

* `go build` to create binary 'lbas(.exe)'.
  
* run `lbas(.exe) -in input/folder/for/images -out output/folder/for/scaled-images -cfg config/file/path/config.json`

* in `config.json`:
        `roi` is focus area in original image.
        `key-area` is mask region in roi area image.
        `valid-range` is detect range ration, top to bottom is from 0.0 to 1.0.
        `thBright` is threshold value for brighter area.
        `thContrast` is threshold value for bright/dark areas. 


