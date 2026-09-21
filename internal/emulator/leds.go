package emulator

import (
	"image/color"
	"sync"
)

// LEDs keeps separate draw and presented buffers like an addressable strip.
type LEDs struct {
	mu     sync.RWMutex
	buffer []color.RGBA
	front  []color.RGBA
}

func NewLEDs(count int) *LEDs {
	if count < 0 {
		panic("LED count must not be negative")
	}
	return &LEDs{buffer: make([]color.RGBA, count), front: make([]color.RGBA, count)}
}

func (l *LEDs) Len() int { return len(l.buffer) }

func (l *LEDs) Set(index int, c color.RGBA) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index >= 0 && index < len(l.buffer) {
		l.buffer[index] = c
	}
}

func (l *LEDs) Display() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	copy(l.front, l.buffer)
	return nil
}

// Color returns the presented color. WS2812 ignores alpha.
func (l *LEDs) Color(index int) color.RGBA {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if index < 0 || index >= len(l.front) {
		return color.RGBA{}
	}
	return l.front[index]
}
