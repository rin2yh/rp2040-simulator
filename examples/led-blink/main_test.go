package main

import (
	"image/color"
	"testing"
)

func TestFill(t *testing.T) {
	leds := make([]color.RGBA, ledCount)
	want := color.RGBA{G: 24}
	fill(leds, want)
	for i, got := range leds {
		if got != want {
			t.Fatalf("LED %d = %v, want %v", i, got, want)
		}
	}
}
