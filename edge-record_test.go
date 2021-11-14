package main

import (
	"fmt"
	"testing"
)

func TestEdgeRecord(t *testing.T) {
	record := NewEdgeRecord("AB", "./in/sample/1.jpg")
	record.AddPtInfo(
		100,
		100,
		100,
		34,
		12,
		44,
	)
	record.AddPtInfo(
		200,
		200,
		200,
		3,
		12,
		44,
	)
	record.Log("./out/AB-edge.json")
}

func TestLoadEdgeRecord1(t *testing.T) {
	edge := LoadLastRecord("./cfg/AB-edge.json")
	fmt.Println(edge)
	fmt.Println(edge.Pts[0].ValAbove)
	fmt.Println(edge.Pts[0].ValBelow)
}
