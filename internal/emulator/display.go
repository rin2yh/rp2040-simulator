// Package emulator implements the zero-kb02 desktop emulator.
package emulator

import (
	"image/color"

	"github.com/rin2yh/rp2040-simulator/internal/bounds"
)

// Display keeps separate draw and presented buffers, just like an OLED.
// Bits use SSD1306 page order: x + (y/8)*width, bit y%8.
type Display struct {
	width, height int16
	frame         frameBuffer[byte]
}

// NewDisplay requires positive dimensions.
func NewDisplay(width, height int16) *Display {
	if width <= 0 || height <= 0 {
		panic("display dimensions must be positive")
	}
	size := int(width) * ((int(height) + 7) / 8)
	return &Display{width: width, height: height, frame: newFrameBuffer[byte](size)}
}

func (d *Display) Size() (int16, int16) { return d.width, d.height }

func (d *Display) SetPixel(x, y int16, c color.RGBA) {
	if !bounds.Between(x, 0, d.width) || !bounds.Between(y, 0, d.height) {
		return
	}
	i := int(x) + int(y/8)*int(d.width)
	mask := byte(1 << uint(y%8))
	d.frame.update(i, func(value *byte) {
		if c.R != 0 || c.G != 0 || c.B != 0 {
			*value |= mask
		} else {
			*value &^= mask
		}
	})
}

func (d *Display) ClearBuffer() { d.frame.clear() }

func (d *Display) Display() error {
	d.frame.display()
	return nil
}

// WriteFrame atomically replaces the drawing and presented buffers.
func (d *Display) WriteFrame(pixels []byte) error {
	return d.frame.write(pixels)
}

// Pixel reads the presented frame, never the uncommitted drawing buffer.
func (d *Display) Pixel(x, y int16) bool {
	if !bounds.Between(x, 0, d.width) || !bounds.Between(y, 0, d.height) {
		return false
	}
	return d.frame.presented(int(x)+int(y/8)*int(d.width))&(1<<uint(y%8)) != 0
}
