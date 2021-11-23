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
		Left     int `json:"left"`
		Top      int `json:"top"`
		Right    int `json:"right"`
		Bottom   int `json:"bottom"`
		AimAreas []struct {
			File       string     `json:"file"`
			Offset     int        `json:"offset"`
			ValidRange [2]float64 `json:"valid-range"`
			ThBright   float64    `json:"thBright"`
			ThContrast float64    `json:"thContrast"`
		} `json:"aim-areas"`
	} `json:"roi"`
}

func (c *Config) String() (s string) {
	s += fmt.Sprintf("%-40s[%v,%v,%v,%v]\n", "aimming area:", c.Roi.Left, c.Roi.Top, c.Roi.Right, c.Roi.Bottom)
	for _, area := range c.Roi.AimAreas {
		s += fmt.Sprintf("  %-38s%v\n", "key area:", area.File)
		s += fmt.Sprintf("  %-38s%v\n", "key area offset:", area.Offset)
		s += fmt.Sprintf("  %-38s%v\n", "checking range:", area.ValidRange)
		s += fmt.Sprintf("  %-38s%v\n", "threshold for brighter area:", area.ThBright)
		s += fmt.Sprintf("  %-38s%v\n\n", "contrast for bright/dark areas:", area.ThContrast)
	}
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
