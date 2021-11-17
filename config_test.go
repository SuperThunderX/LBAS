package main

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig("./cfg/config.json")
	fmt.Println(cfg)
}
