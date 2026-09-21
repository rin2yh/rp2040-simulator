//go:build !tinygo

package board

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Conf2025Badge follows the assembled board in
// https://github.com/tinygo-keeb/workshop-conf2025badge.
// XIAO RP2040; SSD1306 128x64; two SK6812 key LEDs on GPIO0.
func Conf2025Badge() Profile {
	return Profile{
		Name: "conf2025badge", DisplayWidth: 128, DisplayHeight: 64, LEDCount: 2,
		View: View{
			Width: 890, Height: 690,
			Display: image.Rect(371, 185, 551, 275),
			Encoder: image.Pt(693, 209), EncoderRadius: 47,
			Joystick: image.Pt(208, 355), JoystickRadius: 43,
			JoystickLabel: image.Pt(149, 405),
			Boot:          image.Rect(155, 218, 180, 239), Reset: image.Rect(235, 218, 260, 239),
			Keys: []Key{
				{Bounds: image.Rect(635, 300, 751, 416), Code: ebiten.KeyQ, Label: "Q"},
				{Bounds: image.Rect(635, 427, 751, 543), Code: ebiten.KeyW, Label: "W"},
			},
			LEDs:     []image.Point{image.Pt(693, 358), image.Pt(693, 485)},
			DrawBody: drawConf2025Badge,
		},
	}
}

func drawConf2025Badge(dst *ebiten.Image) {
	dst.Fill(rgb(230, 233, 236))
	Label(dst, "conf2025badge", 38, 27, 25, Ink)
	Label(dst, "TinyGo Conference 2025", 39, 61, 14, rgb(89, 96, 104))
	rounded(dst, 123, 145, 665, 440, 18, rgb(182, 190, 196))
	rounded(dst, 120, 135, 665, 440, 15, rgb(30, 45, 47))
	rounded(dst, 622, 143, 146, 417, 13, rgb(36, 101, 172))
	for _, x := range []float32{144, 595} {
		vector.FillCircle(dst, x, 550, 7, rgb(61, 123, 186), true)
	}
	// XIAO module and its USB socket at the top left.
	rounded(dst, 143, 157, 131, 94, 4, rgb(36, 83, 74))
	rounded(dst, 151, 164, 111, 49, 3, rgb(168, 181, 178))
	rounded(dst, 125, 169, 35, 38, 3, rgb(193, 199, 197))
	Label(dst, "XIAO RP2040", 170, 180, 11, Ink)
	for _, x := range []float32{155, 235} {
		rounded(dst, x, 218, 25, 21, 3, rgb(158, 166, 153))
	}
	Label(dst, "BOOT", 149, 256, 10, rgb(205, 219, 210))
	Label(dst, "RESET", 232, 256, 10, rgb(205, 219, 210))
	rounded(dst, 357, 159, 208, 143, 4, rgb(21, 31, 37))
	Label(dst, "GND  VDD  SCK  SDA", 378, 166, 10, rgb(185, 203, 198))
	Label(dst, "TinyGo Conference 2025", 300, 313, 22, rgb(209, 226, 215))
	rounded(dst, 286, 356, 321, 191, 13, rgb(215, 231, 220))
	// The buzzer is shown as a physical part; audio is not simulated.
	vector.FillCircle(dst, 208, 477, 39, rgb(13, 21, 24), true)
	vector.FillCircle(dst, 208, 470, 34, rgb(47, 59, 62), true)
	vector.FillCircle(dst, 208, 465, 5, rgb(9, 15, 18), true)
	Label(dst, "BUZZER", 181, 520, 10, rgb(191, 208, 200))
	drawKeycap(dst, 635, 300)
	drawKeycap(dst, 635, 427)
	Label(dst, "Keys: Q / W    Stick: arrows    Knob: [ ]    Reset: F5", 126, 648, 13, rgb(99, 105, 111))
}
