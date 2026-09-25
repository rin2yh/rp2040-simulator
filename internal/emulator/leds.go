package emulator

import (
	"image/color"
)

// LEDs keeps separate draw and presented buffers like an addressable strip.
type LEDs struct {
	frame frameBuffer[color.RGBA]
}

func NewLEDs(count int) *LEDs {
	if count < 0 {
		panic("LED count must not be negative")
	}
	return &LEDs{frame: newFrameBuffer[color.RGBA](count)}
}

func (l *LEDs) Len() int { return l.frame.len() }

func (l *LEDs) Set(index int, c color.RGBA) {
	if index >= 0 && index < l.Len() {
		l.frame.update(index, func(value *color.RGBA) { *value = c })
	}
}

func (l *LEDs) Display() error {
	l.frame.display()
	return nil
}

func (l *LEDs) WriteColors(colors []color.RGBA) error { return l.frame.write(colors) }

// Color returns the presented color. WS2812 ignores alpha.
func (l *LEDs) Color(index int) color.RGBA {
	if index < 0 || index >= l.Len() {
		return color.RGBA{}
	}
	return l.frame.presented(index)
}
