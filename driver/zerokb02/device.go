// Package zerokb02 exposes the zero-kb02 controls and OLED on desktop and TinyGo.
// LED output uses the existing driver/ws2812 package with machine.GPIO1.
package zerokb02

import "image/color"

const Width, Height = 128, 64

type State struct {
	Keys                 [12]bool // row-major: QWER, ASDF, ZXCV
	EncoderDelta         int      // detents since this Device's previous Read
	EncoderPressed       bool
	JoystickX, JoystickY float32 // -1 to 1; up is negative Y
	JoystickPressed      bool
}

// Device owns one control reader and one OLED framebuffer.
type Device struct {
	frame       [Width * Height / 8]byte
	position    int
	generation  uint64
	initialized bool
	hardware    hardware
}

func (d *Device) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return
	}
	i, bit := int(x)+int(y/8)*Width, byte(1<<uint(y%8))
	if c.R != 0 || c.G != 0 || c.B != 0 {
		d.frame[i] |= bit
	} else {
		d.frame[i] &^= bit
	}
}

func (d *Device) ClearBuffer() { clear(d.frame[:]) }
