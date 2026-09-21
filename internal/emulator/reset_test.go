//go:build !tinygo

package emulator

import (
	"image"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rin2yh/rp2040-simulator/internal/board"
)

func TestResetClearsDevicesAndTogglesBootloader(t *testing.T) {
	profile := board.Profile{DisplayWidth: 128, DisplayHeight: 64, LEDCount: 1, View: board.View{
		Width: 400, Height: 300, Display: image.Rect(0, 0, 128, 64), EncoderRadius: 30,
		Keys: []board.Key{{Bounds: image.Rect(150, 150, 190, 190)}}, LEDs: []image.Point{image.Pt(170, 170)}, DrawBody: func(*ebiten.Image) {},
	}}
	g, err := newGame(profile)
	if err != nil {
		t.Fatal(err)
	}
	g.encoder.Rotate(5)
	g.buttons[0].Down = true
	g.joystick.X = 1
	g.display.SetPixel(0, 0, color.RGBA{R: 255})
	g.display.Display()
	g.leds.Set(0, color.RGBA{R: 255})
	g.leds.Display()
	if err := g.restart(false); err != nil {
		t.Fatal(err)
	}
	if g.encoder.Position() != 0 || g.buttons[0].Pressed() || g.joystick.X != 0 || g.display.Pixel(0, 0) || g.leds.Color(0).R != 0 {
		t.Fatal("device state survived reset")
	}
	if err := g.restart(true); err != nil {
		t.Fatal(err)
	}
	if !g.bootloader {
		t.Fatal("bootloader did not start")
	}
	if err := g.restart(false); err != nil {
		t.Fatal(err)
	}
	if g.bootloader {
		t.Fatal("cannot exit bootloader with reset")
	}
}
