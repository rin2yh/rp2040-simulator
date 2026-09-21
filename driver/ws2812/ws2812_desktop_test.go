//go:build !tinygo

package ws2812

import (
	"image/color"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
	"github.com/rin2yh/rp2040-simulator/machine"
)

func TestWriteColorsUsesDeviceRPC(t *testing.T) {
	previous := call
	defer func() { call = previous }()

	want := []color.RGBA{{R: 1}, {G: 2}, {B: 3}}
	call = func(method string, args, _ any) error {
		if method != "LEDs.WriteColors" {
			t.Fatalf("method = %q", method)
		}
		got := args.(*bridge.WriteLEDsArgs).Colors
		if len(got) != len(want) {
			t.Fatalf("color count = %d, want %d", len(got), len(want))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("color %d = %v, want %v", i, got[i], want[i])
			}
		}
		return nil
	}

	device := NewWS2812(machine.GPIO1)
	if err := device.WriteColors(want); err != nil {
		t.Fatal(err)
	}
}

func TestWriteColorsRejectsUnknownPin(t *testing.T) {
	if err := NewWS2812(machine.Pin(2)).WriteColors(nil); err == nil {
		t.Fatal("unknown pin was accepted")
	}
}
