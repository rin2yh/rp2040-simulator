//go:build !tinygo

// Package board contains the dimensions and desktop drawing for each board.
package board

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Profile struct {
	Name                        string
	DisplayWidth, DisplayHeight int16
	LEDCount                    int
	View                        View
}

type View struct {
	Width, Height  int
	Display        image.Rectangle
	Encoder        image.Point
	EncoderRadius  int
	Keys           []Key
	LEDs           []image.Point
	Joystick       image.Point
	JoystickRadius int
	Boot, Reset    image.Rectangle
	DrawBody       func(*ebiten.Image)
}

type Key struct {
	Bounds image.Rectangle
	Code   ebiten.Key
	Label  string
}
