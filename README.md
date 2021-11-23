# LBAS

Liquid Boundary Analysis System

Compute & Mark 1st liquid border, computed scale value (0.0-1.0) is added on output image file name.

* `go build` to create binary 'lbas(.exe)'.

* for detect area config creating:
   
   1. prepare pure color (R, G, B, Y, C, M, K, W) marked lines on model image. An example is './in/sample/marked/sample.jpg'.
   2. run `lbas(.exe) -ca true -caimg path/to/your/lines-image.jpg -caout path/to/output.json -caclr R/G/B...`
   3. if multiple-lines in model image, run above command line with correct arguments multiple times to get different area.json.
   4. open config.json, in "aim-areas", add a new array element. Its "file" is the output json path created by previous step. other arguments are this line area detecting arguments. 
  
* run `lbas(.exe) -in input/folder/for/images -out output/folder/for/scaled-images -cfg config/file/path/config.json`

* in `config.json`:
        `roi` is focus area in original image. if all are 0, take original image size to process.
        `aim-areas` are lines mask region in roi image.
        `offset` is adjusting x position for its line area. negative is moving to left, positive is moving to right.
        `valid-range` is detect range ration, top to bottom is from 0.0 to 1.0.
        `thBright` is threshold value for brighter area.
        `thContrast` is threshold value for bright/dark areas. 


