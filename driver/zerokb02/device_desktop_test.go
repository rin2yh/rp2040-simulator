//go:build !tinygo

package zerokb02

import (
	"image/color"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

func TestReadAndDisplay(t *testing.T) {
	old := call
	defer func() { call = old }()
	position, generation := 3, uint64(1)
	call = func(method string, args, reply any) error {
		switch method {
		case "Inputs.Read":
			*reply.(*bridge.Inputs) = bridge.Inputs{Keys: [12]bool{true}, EncoderPosition: position, Generation: generation}
		case "Display.WriteFrame":
			frame := args.(*bridge.DisplayFrame).Pixels
			if len(frame) != 1024 || frame[127+7*128] != 0x80 {
				t.Fatalf("frame bytes invalid: %d", len(frame))
			}
		default:
			t.Fatalf("unexpected method %q", method)
		}
		return nil
	}
	d, _ := New()
	first, err := d.Read()
	if err != nil || first.EncoderDelta != 3 || !first.Keys[0] {
		t.Fatalf("first read = %+v, %v", first, err)
	}
	position = 5
	second, err := d.Read()
	if err != nil || second.EncoderDelta != 2 {
		t.Fatalf("second read = %+v, %v", second, err)
	}
	generation++
	position = 0
	reset, err := d.Read()
	if err != nil || reset.EncoderDelta != 0 {
		t.Fatalf("reset read = %+v, %v", reset, err)
	}
	d.SetPixel(127, 63, color.RGBA{R: 1})
	if err := d.Display(); err != nil {
		t.Fatal(err)
	}
}
