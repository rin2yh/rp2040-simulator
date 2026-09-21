//go:build !tinygo

// Package board contains the dimensions and desktop drawing for each board.
package board

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const Environment = "RP2040_SIMULATOR_BOARD"

func Lookup(name string) (Profile, error) {
	switch name {
	case "", "zero-kb02":
		return ZeroKB02(), nil
	case "conf2025badge":
		return Conf2025Badge(), nil
	default:
		return Profile{}, fmt.Errorf("unknown board %q (want zero-kb02 or conf2025badge)", name)
	}
}

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
	JoystickLabel  image.Point
	JoystickRadius int
	Boot, Reset    image.Rectangle
	DrawBody       func(*ebiten.Image)
}

type Key struct {
	Bounds image.Rectangle
	Code   ebiten.Key
	Label  string
}
