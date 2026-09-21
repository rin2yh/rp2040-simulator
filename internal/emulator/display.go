// Package emulator implements the zero-kb02 desktop emulator.
package emulator

import "image/color"

// Display keeps separate draw and presented buffers, just like an OLED.
// Bits use SSD1306 page order: x + (y/8)*width, bit y%8.
type Display struct {
	width, height int16
	buffer, front []byte
}

// NewDisplay requires positive dimensions.
func NewDisplay(width, height int16) *Display {
	if width <= 0 || height <= 0 {
		panic("display dimensions must be positive")
	}
	size := int(width) * ((int(height) + 7) / 8)
	return &Display{width: width, height: height, buffer: make([]byte, size), front: make([]byte, size)}
}

func (d *Display) Size() (int16, int16) { return d.width, d.height }

func (d *Display) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || y < 0 || x >= d.width || y >= d.height {
		return
	}
	i := int(x) + int(y/8)*int(d.width)
	mask := byte(1 << uint(y%8))
	if c.R != 0 || c.G != 0 || c.B != 0 {
		d.buffer[i] |= mask
	} else {
		d.buffer[i] &^= mask
	}
}

func (d *Display) ClearBuffer() { clear(d.buffer) }

func (d *Display) Display() error {
	copy(d.front, d.buffer)
	return nil
}

// Pixel reads the presented frame, never the uncommitted drawing buffer.
func (d *Display) Pixel(x, y int16) bool {
	if x < 0 || y < 0 || x >= d.width || y >= d.height {
		return false
	}
	return d.front[int(x)+int(y/8)*int(d.width)]&(1<<uint(y%8)) != 0
}
