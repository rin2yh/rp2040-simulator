package emulator

import (
	"image/color"
	"testing"

	"tinygo.org/x/drivers/ssd1306"
)

// This transport is only a test double for comparing the real driver's
// framebuffer. The emulator has no I2C or SSD1306 command implementation.
type discardI2C struct{}

func (discardI2C) Tx(uint16, []byte, []byte) error { return nil }

func TestSSD1306Parity(t *testing.T) {
	real := ssd1306.NewI2C(discardI2C{})
	real.Configure(ssd1306.Config{Width: 128, Height: 64, Address: 0x3c})
	emulated := NewDisplay(128, 64)
	colors := []color.RGBA{{}, {R: 1}, {G: 1}, {B: 1}, {A: 255}, {R: 255, A: 255}}
	// Include out-of-bounds coordinates, all page boundaries, and opaque
	// black/transparent non-black, whose alpha is ignored by the real driver.
	for y := int16(-1); y <= 64; y++ {
		for x := int16(-1); x <= 128; x++ {
			c := colors[(int(x)+int(y)+6)%len(colors)]
			real.SetPixel(x, y, c)
			emulated.SetPixel(x, y, c)
		}
	}
	if err := emulated.Display(); err != nil {
		t.Fatal(err)
	}
	for y := int16(-1); y <= 64; y++ {
		for x := int16(-1); x <= 128; x++ {
			if got, want := emulated.Pixel(x, y), real.GetPixel(x, y); got != want {
				t.Fatalf("pixel (%d,%d): got %t, want %t", x, y, got, want)
			}
		}
	}
}

func TestDisplayCommitsOnlyOnFlush(t *testing.T) {
	d := NewDisplay(128, 64)
	d.SetPixel(127, 63, color.RGBA{R: 255})
	if d.Pixel(127, 63) {
		t.Fatal("uncommitted pixel became visible")
	}
	if err := d.Display(); err != nil {
		t.Fatal(err)
	}
	if !d.Pixel(127, 63) {
		t.Fatal("committed pixel is missing")
	}
	d.ClearBuffer()
	if !d.Pixel(127, 63) {
		t.Fatal("clear modified the presented frame")
	}
	if err := d.Display(); err != nil {
		t.Fatal(err)
	}
	if d.Pixel(127, 63) {
		t.Fatal("clear was not committed")
	}
}
