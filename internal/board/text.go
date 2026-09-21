//go:build !tinygo

package board

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

var Ink = color.RGBA{R: 56, G: 62, B: 68, A: 255}
var fontSource = func() *text.GoTextFaceSource {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	} // The font is embedded, not user-supplied.
	return source
}()

// Label draws desktop UI text. It is not used for the emulated OLED, which
// always displays the application's original monochrome pixels.
func Label(dst *ebiten.Image, label string, x, y, size float64, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(dst, label, &text.GoTextFace{Source: fontSource, Size: size}, op)
}
