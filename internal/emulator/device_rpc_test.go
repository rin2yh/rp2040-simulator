//go:build !tinygo

package emulator

import (
	"image/color"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

func TestDeviceServicesUseGameState(t *testing.T) {
	g := &game{profile: board.ZeroKB02(), buttons: make([]Button, 12), display: NewDisplay(128, 64)}
	g.buttons[7].Down = true
	g.encoder.Rotate(3)
	g.encoder.SetPressed(true)
	g.joystick = Joystick{X: -0.5, Y: 1, Down: true}
	var got bridge.Inputs
	if err := (&inputService{game: g}).Read(&struct{}{}, &got); err != nil {
		t.Fatal(err)
	}
	if !got.Keys[7] || got.Keys[6] || got.EncoderPosition != 3 || !got.EncoderPressed || got.JoystickX != -0.5 || got.JoystickY != 1 || !got.JoystickPressed {
		t.Fatalf("inputs = %+v", got)
	}
	frame := make([]byte, 1024)
	frame[127+7*128] = 0x80
	service := displayService{display: g.display}
	if err := service.WriteFrame(&bridge.DisplayFrame{Pixels: frame}, &struct{}{}); err != nil {
		t.Fatal(err)
	}
	if !g.display.Pixel(127, 63) || g.display.Pixel(127, 62) {
		t.Fatal("display did not present the page-order frame")
	}
	if err := service.WriteFrame(&bridge.DisplayFrame{Pixels: frame[:1]}, &struct{}{}); err == nil {
		t.Fatal("short frame accepted")
	}
	if !g.display.Pixel(127, 63) {
		t.Fatal("invalid frame modified display")
	}
	g.display.SetPixel(0, 0, color.RGBA{R: 255})
	if g.display.Pixel(0, 0) {
		t.Fatal("uncommitted pixel is visible")
	}
}
