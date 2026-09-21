//go:build !tinygo

package board

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Top view based on the zero-kb02 assembly photos and KiCad layout.
// PCB: 93 x 93 mm. Key pitch: 19.05 mm. Render scale: 6.8 pixels/mm.
// Electrical pins and these desktop hit regions both belong to this profile.
func ZeroKB02() Profile {
	view := View{
		Width: 890, Height: 860,
		Display: image.Rect(440, 197, 606, 280),
		Encoder: image.Pt(684, 213), EncoderRadius: 47,
		Joystick: image.Pt(215, 211), JoystickRadius: 43,
		Boot: image.Rect(329, 247, 359, 279), Reset: image.Rect(386, 247, 416, 279),
		DrawBody: drawBody,
	}
	codes := []ebiten.Key{ebiten.KeyQ, ebiten.KeyW, ebiten.KeyE, ebiten.KeyR, ebiten.KeyA, ebiten.KeyS, ebiten.KeyD, ebiten.KeyF, ebiten.KeyZ, ebiten.KeyX, ebiten.KeyC, ebiten.KeyV}
	for i, label := range []string{"Q", "W", "E", "R", "A", "S", "D", "F", "Z", "X", "C", "V"} {
		bounds := keyBounds(i)
		view.Keys = append(view.Keys, Key{Bounds: bounds, Code: codes[i], Label: label})
		view.LEDs = append(view.LEDs, image.Pt(bounds.Min.X+58, bounds.Min.Y+58))
	}
	return Profile{Name: "zero-kb02", DisplayWidth: 128, DisplayHeight: 64, LEDCount: 12, View: view}
}

func keyBounds(i int) image.Rectangle {
	x := int(126+(116.6813-98.72+float64(i%4)*19.05)*6.8) - 58
	y := int(115+(95.25-52.46+float64(i/4)*19.05)*6.8) - 58
	return image.Rect(x, y, x+116, y+116)
}

func rgb(r, g, b uint8) color.RGBA { return color.RGBA{R: r, G: g, B: b, A: 255} }

func rounded(dst *ebiten.Image, x, y, w, h, radius float32, c color.Color) {
	vector.FillRect(dst, x+radius, y, w-2*radius, h, c, false)
	vector.FillRect(dst, x, y+radius, w, h-2*radius, c, false)
	for _, center := range [][2]float32{{x + radius, y + radius}, {x + w - radius, y + radius}, {x + radius, y + h - radius}, {x + w - radius, y + h - radius}} {
		vector.FillCircle(dst, center[0], center[1], radius, c, true)
	}
}

func drawBody(dst *ebiten.Image) {
	dst.Fill(rgb(230, 233, 236))
	Label(dst, "zero-kb02", 38, 27, 25, rgb(34, 40, 45))
	Label(dst, "Interactive device emulator", 39, 61, 14, rgb(89, 96, 104))
	// A short cable and exposed RP2040-Zero, as on the assembled board.
	rounded(dst, 357, 89, 26, 91, 9, rgb(58, 61, 64))
	// Layered case edge and subtle contact shadow.
	rounded(dst, 127, 132, 642, 632, 20, rgb(194, 199, 204))
	rounded(dst, 126, 122, 632, 634, 17, rgb(22, 24, 26))
	rounded(dst, 126, 115, 632, 630, 16, rgb(43, 45, 48))
	rounded(dst, 131, 119, 622, 620, 13, rgb(37, 39, 42))
	// Printed case texture, kept subtle enough that the OLED stays legible.
	for y := float32(123); y < 735; y += 4 {
		vector.StrokeLine(dst, 142, y, 742, y, 0.5, rgb(44, 46, 48), false)
	}
	// The raised joystick enclosure occupies the top-left of the case.
	rounded(dst, 145, 135, 136, 181, 5, rgb(18, 20, 22))
	rounded(dst, 148, 137, 130, 170, 4, rgb(43, 45, 48))
	vector.FillCircle(dst, 215, 226, 43, rgb(13, 15, 17), true)
	// RP2040-Zero board, header pads, USB-C shell and the two small switches.
	rounded(dst, 316, 139, 109, 171, 3, rgb(12, 16, 18))
	rounded(dst, 320, 135, 101, 168, 3, rgb(24, 72, 87))
	for y := float32(146); y < 297; y += 17 {
		for _, x := range []float32{324, 417} {
			vector.FillCircle(dst, x, y, 3, rgb(190, 172, 112), true)
			vector.FillCircle(dst, x, y, 1, rgb(17, 25, 27), true)
		}
	}
	for x := float32(337); x < 413; x += 17 {
		vector.FillCircle(dst, x, 296, 3, rgb(190, 172, 112), true)
	}
	rounded(dst, 349, 126, 43, 47, 4, rgb(174, 180, 178))
	rounded(dst, 352, 124, 37, 10, 4, rgb(83, 90, 91))
	rounded(dst, 357, 128, 27, 4, 2, rgb(22, 28, 30))
	rounded(dst, 353, 199, 34, 34, 2, rgb(23, 28, 30))
	for _, x := range []float32{337, 394} {
		rounded(dst, x, 254, 13, 18, 2, rgb(161, 162, 152))
		rounded(dst, x+3, 257, 7, 12, 2, rgb(219, 216, 180))
	}
	Label(dst, "RP2040", 348, 183, 9, rgb(202, 213, 202))
	Label(dst, "ZERO", 354, 236, 10, rgb(202, 213, 202))
	Label(dst, "BOOT", 327, 280, 9, rgb(202, 213, 202))
	Label(dst, "RESET", 386, 280, 9, rgb(202, 213, 202))
	// OLED opening; live pixels are overlaid by the generic emulator runner.
	rounded(dst, 429, 182, 188, 112, 4, rgb(13, 15, 19))
	rounded(dst, 434, 188, 178, 100, 3, rgb(19, 26, 35))
	// 4 x 3 keys. Centers are taken from the KiCad SW1..SW12 positions.
	for row := range 3 {
		for col := range 4 {
			r := keyBounds(row*4 + col)
			drawKeycap(dst, float32(r.Min.X), float32(r.Min.Y))
		}
	}
	Label(dst, "Keys: QWER / ASDF / ZXCV    Stick: arrows    Knob: [ ]    Reset: F5", 126, 813, 13, rgb(99, 105, 111))
}

func drawKeycap(dst *ebiten.Image, x, y float32) {
	// Translucent blank keycaps reveal the switch stem and RGB LED beneath it.
	rounded(dst, x-3, y+3, 122, 122, 9, rgb(15, 18, 20))
	rounded(dst, x, y, 116, 116, 9, rgb(114, 128, 133))
	rounded(dst, x+2, y+1, 112, 108, 8, rgb(177, 188, 190))
	rounded(dst, x+6, y+5, 104, 99, 7, rgb(133, 151, 154))
	rounded(dst, x+17, y+13, 82, 77, 9, rgb(147, 166, 167))
	vector.FillCircle(dst, x+58, y+48, 22, rgb(125, 147, 150), true)
	vector.FillCircle(dst, x+58, y+48, 17, rgb(160, 183, 184), true)
	vector.FillRect(dst, x+54, y+35, 8, 26, rgb(100, 134, 141), false)
	vector.FillRect(dst, x+45, y+44, 26, 8, rgb(100, 134, 141), false)
	vector.StrokeLine(dst, x+10, y+8, x+10, y+98, 2, rgb(194, 205, 204), true)
	vector.StrokeLine(dst, x+106, y+8, x+106, y+98, 2, rgb(165, 183, 184), true)
	vector.StrokeLine(dst, x+20, y+15, x+95, y+15, 2, rgb(182, 199, 197), true)
	rounded(dst, x+22, y+87, 72, 8, 3, rgb(100, 129, 136))
	vector.StrokeLine(dst, x+11, y+106, x+105, y+106, 2, rgb(196, 207, 205), true)
}
