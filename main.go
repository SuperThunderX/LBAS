package main

import "fmt"

func main() {

	cfg := LoadConfig("./cfg/config.json")
	fmt.Println(cfg)

}
