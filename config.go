package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"os"
)

type Config struct {
	Roi struct {
		Left   int `json:"left"`
		Top    int `json:"top"`
		Right  int `json:"right"`
		Bottom int `json:"bottom"`
	} `json:"roi"`
	KeyArea       string     `json:"key-area"`
	KeyAreaOffset int        `json:"key-area-offset"`
	ValidRange    [2]float64 `json:"valid-range"`
	ThBright      float64    `json:"thBright"`
	ThContrast    float64    `json:"thContrast"`
}

func (c *Config) String() (s string) {
	s += fmt.Sprintf("%-40s%v\n", "checking area:", c.Roi)
	s += fmt.Sprintf("%-40s%v\n", "key area:", c.KeyArea)
	s += fmt.Sprintf("%-40s%v\n", "key area offset:", c.KeyAreaOffset)
	s += fmt.Sprintf("%-40s%v\n", "checking range:", c.ValidRange)
	s += fmt.Sprintf("%-40s%v\n", "threshold for brighter area:", c.ThBright)
	s += fmt.Sprintf("%-40s%v\n", "contrast for bright/dark areas:", c.ThContrast)
	return s
}

func (c *Config) ROIRect() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: c.Roi.Left,
			Y: c.Roi.Top,
		},
		Max: image.Point{
			X: c.Roi.Right,
			Y: c.Roi.Bottom,
		},
	}
}

func LoadConfig(cfgfile string) *Config {
	data, err := os.ReadFile(cfgfile)
	if err != nil {
		log.Fatalln(err)
	}
	cfg := new(Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		log.Fatalln(err)
	}
	return cfg
}
