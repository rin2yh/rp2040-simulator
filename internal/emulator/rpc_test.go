//go:build !tinygo

package emulator

import (
	"image/color"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

func TestLEDServiceCommitsCompleteFrame(t *testing.T) {
	leds := NewLEDs(2)
	service := ledService{leds: leds}
	want := []color.RGBA{{R: 12}, {B: 34}}
	if err := service.WriteColors(&bridge.WriteLEDsArgs{Colors: want}, &struct{}{}); err != nil {
		t.Fatal(err)
	}
	for i, c := range want {
		if got := leds.Color(i); got != c {
			t.Fatalf("LED %d = %v, want %v", i, got, c)
		}
	}
	if err := service.WriteColors(&bridge.WriteLEDsArgs{Colors: want[:1]}, &struct{}{}); err == nil {
		t.Fatal("wrong LED count was accepted")
	}
}
