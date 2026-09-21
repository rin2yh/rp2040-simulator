package emulator

import (
	"image/color"
	"testing"
)

func TestLEDsCommitCompleteFrame(t *testing.T) {
	leds := NewLEDs(2)
	red := color.RGBA{R: 32}
	leds.Set(0, red)
	leds.Set(-1, color.RGBA{G: 32})
	leds.Set(2, color.RGBA{B: 32})
	if got := leds.Color(0); got != (color.RGBA{}) {
		t.Fatalf("uncommitted LED became visible: %v", got)
	}
	if err := leds.Display(); err != nil {
		t.Fatal(err)
	}
	if got := leds.Color(0); got != red {
		t.Fatalf("committed LED = %v, want %v", got, red)
	}
	if leds.Len() != 2 || leds.Color(-1) != (color.RGBA{}) || leds.Color(2) != (color.RGBA{}) {
		t.Fatal("LED bounds are incorrect")
	}
}
