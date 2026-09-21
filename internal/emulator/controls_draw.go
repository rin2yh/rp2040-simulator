//go:build !tinygo

package emulator

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rin2yh/rp2040-simulator/internal/board"
)

const defaultHint = "Click keys; drag the knob or joystick. Hover a control for help."

func (g *game) updateHint(point image.Point) {
	g.hint = defaultHint
	hover := false
	switch {
	case point.In(g.view.Boot):
		g.hint = "BOOT: click to arm, then click RESET. Keyboard: hold B and press F5."
		hover = true
	case point.In(g.view.Reset):
		g.hint = "RESET: restart the application from scratch. Keyboard: F5."
		hover = true
	case point.In(g.view.Display):
		g.hint = "OLED: click to enlarge or press O. Click again or Esc to return."
		hover = true
	case InsideCircle(point, g.view.Encoder, g.view.EncoderRadius):
		g.hint = "Knob: drag to turn / scroll here. Click to push; right-click or Space to hold."
		hover = true
	case InsideCircle(point, g.view.Joystick, g.view.JoystickRadius):
		g.hint = "Stick: drag / arrow keys. Click to push; right-click or Shift to hold."
		hover = true
	default:
		for i, key := range g.view.Keys {
			if point.In(key.Bounds) {
				g.hint = fmt.Sprintf("Key %d: hold the mouse button or %s. The OLED shows the input.", i+1, key.Label)
				hover = true
				break
			}
		}
	}
	if g.bootArmed || g.bootHeld {
		g.hint = "BOOT held. Press RESET to enter the simulated bootloader; click BOOT to disarm."
	}
	if hover && !g.zoom {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	}
}

func (g *game) drawInputs(dst *ebiten.Image) {
	for i, point := range g.view.LEDs {
		c := g.leds.Color(i)
		if c.R == 0 && c.G == 0 && c.B == 0 {
			continue
		}
		peak := max(c.R, max(c.G, c.B))
		scale := float32(255) / float32(peak)
		light := color.RGBA{R: uint8(float32(c.R) * scale), G: uint8(float32(c.G) * scale), B: uint8(float32(c.B) * scale), A: 255}
		glow := color.RGBA{R: light.R, G: light.G, B: light.B, A: 72}
		vector.FillCircle(dst, float32(point.X), float32(point.Y), 22, glow, true)
		vector.FillCircle(dst, float32(point.X), float32(point.Y), 9, light, true)
	}
	for i, key := range g.view.Keys {
		r := key.Bounds
		c := board.Ink
		if g.buttons[i].Down {
			vector.FillRect(dst, float32(r.Min.X+3), float32(r.Min.Y+3), float32(r.Dx()-6), float32(r.Dy()-6), color.RGBA{G: 44, B: 55, A: 125}, false)
			vector.StrokeRect(dst, float32(r.Min.X+3), float32(r.Min.Y+3), float32(r.Dx()-6), float32(r.Dy()-6), 3, color.RGBA{R: 108, G: 228, B: 232, A: 255}, true)
			c = color.RGBA{R: 240, G: 255, B: 255, A: 255}
		}
		board.Label(dst, key.Label, float64(r.Min.X+13), float64(r.Min.Y+77), 16, c)
	}
	if g.view.JoystickRadius > 0 {
		x := float32(g.view.Joystick.X) + g.joystick.X*24
		y := float32(g.view.Joystick.Y) + g.joystick.Y*24
		c := color.RGBA{R: 199, G: 201, B: 187, A: 255}
		if g.joystick.Down {
			y += 4
			c = color.RGBA{R: 147, G: 179, B: 172, A: 255}
		}
		vector.FillCircle(dst, x, y+9, 36, color.RGBA{R: 104, G: 115, B: 109, A: 255}, true)
		vector.FillCircle(dst, x, y, 36, color.RGBA{R: 183, G: 188, B: 173, A: 255}, true)
		vector.FillCircle(dst, x, y-1, 30, c, true)
		board.Label(dst, fmt.Sprintf("X %+.1f  Y %+.1f", g.joystick.X, g.joystick.Y), float64(g.view.JoystickLabel.X), float64(g.view.JoystickLabel.Y), 12, color.RGBA{R: 185, G: 196, B: 199, A: 255})
	}
	for _, button := range []struct {
		rect image.Rectangle
		on   bool
	}{
		{g.view.Boot, g.bootArmed || g.bootHeld}, {g.view.Reset, g.resetFlash > 0},
	} {
		if button.on && !button.rect.Empty() {
			r := button.rect
			vector.StrokeRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), 2, color.RGBA{R: 255, G: 200, B: 82, A: 255}, true)
		}
	}
	board.Label(dst, g.hint, 126, float64(g.view.Height-73), 14, board.Ink)
}
