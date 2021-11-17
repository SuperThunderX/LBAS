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
	KeyArea string `json:"key-area"`
}

func (c *Config) String() (s string) {
	s += fmt.Sprintln(c.Roi)
	s += fmt.Sprintln(c.KeyArea)
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
