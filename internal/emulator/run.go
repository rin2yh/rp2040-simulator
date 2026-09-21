//go:build !tinygo

// Package emulator displays logical device state using Ebitengine.
package emulator

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rin2yh/rp2040-simulator/internal/board"
)

type game struct {
	profile               board.Profile
	view                  board.View
	display               *Display
	leds                  *LEDs
	buzzer                *Buzzer
	encoder               Encoder
	oled                  *ebiten.Image
	pixels                []byte
	zoom                  bool
	buttons               []Button
	joystick              Joystick
	pointer               *Pointer
	bootArmed, bootloader bool
	bootHeld              bool
	resetFlash            int
	hint                  string
	testUpdate            func(*game) error
}

// Run owns the main thread and accepts device-level updates over local RPC.
func Run(profile board.Profile) error {
	g, err := newGame(profile)
	if err != nil {
		return err
	}
	return run(g)
}

func run(g *game) error {
	listener, err := startRPC(g.leds, g.buzzer)
	if err != nil {
		return fmt.Errorf("start emulator RPC server: %w", err)
	}
	defer listener.Close()
	defer g.buzzer.closeAudio()
	w, h := g.Layout(0, 0)
	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle(g.profile.Name + " | TinyGo device emulator")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetTPS(60)
	return ebiten.RunGame(g)
}

func newGame(profile board.Profile) (*game, error) {
	view := profile.View
	if profile.DisplayWidth <= 0 || profile.DisplayHeight <= 0 || profile.LEDCount < 0 {
		return nil, errors.New("invalid board profile")
	}
	if view.Width <= 0 || view.Height <= 0 || view.Display.Empty() || view.EncoderRadius <= 0 || view.DrawBody == nil {
		return nil, errors.New("invalid board view")
	}
	g := &game{profile: profile, view: view,
		display: NewDisplay(profile.DisplayWidth, profile.DisplayHeight),
		pixels:  make([]byte, int(profile.DisplayWidth)*int(profile.DisplayHeight)*4)}
	if len(view.LEDs) != profile.LEDCount {
		return nil, errors.New("LED layout does not match board profile")
	}
	if profile.LEDCount > 0 {
		g.leds = NewLEDs(profile.LEDCount)
	}
	if profile.HasBuzzer {
		g.buzzer = &Buzzer{}
	}
	g.buttons = make([]Button, len(view.Keys))
	bounds := make([]image.Rectangle, len(g.buttons))
	for i := range g.buttons {
		bounds[i] = view.Keys[i].Bounds
	}
	g.pointer = NewPointer(PointerLayout{Keys: bounds, Encoder: view.Encoder,
		EncoderRadius: view.EncoderRadius, Joystick: view.Joystick, JoystickRadius: view.JoystickRadius})
	return g, g.restart(false)
}

func (g *game) restart(bootloader bool) error {
	g.buzzer.reset(bootloader)
	g.bootloader, g.bootArmed, g.zoom = bootloader, false, false
	g.encoder = Encoder{}
	g.joystick = Joystick{}
	g.display.ClearBuffer()
	if err := g.display.Display(); err != nil {
		return err
	}
	if g.leds != nil {
		for i := 0; i < g.leds.Len(); i++ {
			g.leds.Set(i, color.RGBA{})
		}
		if err := g.leds.Display(); err != nil {
			return err
		}
	}
	clear(g.buttons)
	g.pointer.Reset()
	return nil
}

func (g *game) Update() error {
	if err := g.buzzer.updateAudio(); err != nil {
		return err
	}
	if g.resetFlash > 0 {
		g.resetFlash--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if g.zoom {
			g.zoom = false
			return nil
		}
		return ebiten.Termination
	}
	mx, my := ebiten.CursorPosition()
	point := image.Pt(mx, my)
	focused := ebiten.IsFocused()
	g.bootHeld = focused && ebiten.IsKeyPressed(ebiten.KeyB)
	click := focused && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	if !g.zoom && click && point.In(g.view.Boot) {
		g.bootArmed = !g.bootArmed
	}
	if focused && (inpututil.IsKeyJustPressed(ebiten.KeyF5) || (!g.zoom && click && point.In(g.view.Reset))) {
		if err := g.restart(g.bootArmed || g.bootHeld); err != nil {
			return err
		}
		g.resetFlash = 12
		g.pointer.Step(PointerInput{Down: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)})
		return nil
	}
	toggleZoom := !g.bootloader && focused && (inpututil.IsKeyJustPressed(ebiten.KeyO) ||
		(click && (g.zoom || point.In(g.view.Display))))
	if toggleZoom {
		g.zoom = !g.zoom
	}
	_, wheel := ebiten.Wheel()
	enabled := focused && !g.zoom && !toggleZoom && !g.bootloader
	g.pointer.Step(PointerInput{Position: point, Down: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		Right: ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight), Wheel: wheel, Enabled: enabled})
	g.updateDevices(enabled)
	g.updateHint(point)
	if g.bootloader {
		return nil
	}
	if g.testUpdate != nil {
		return g.testUpdate(g)
	}
	return nil
}

// Merge pointer and keyboard state after handling reset, focus and zoom.
func (g *game) updateDevices(enabled bool) {
	g.encoder.Rotate(g.pointer.Delta)
	g.encoder.SetPressed(g.pointer.EncoderDown || (enabled && (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter))))
	for i, key := range g.view.Keys {
		g.buttons[i].Down = g.pointer.Keys[i] || (enabled && ebiten.IsKeyPressed(key.Code))
	}
	g.joystick.X, g.joystick.Y, g.joystick.Down = g.pointer.X, g.pointer.Y, g.pointer.JoystickDown
	if enabled {
		if inpututil.IsKeyJustPressed(ebiten.KeyBracketLeft) {
			g.encoder.Rotate(-1)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBracketRight) {
			g.encoder.Rotate(1)
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.joystick.X -= 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.joystick.X += 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.joystick.Y -= 1
		}
		if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.joystick.Y += 1
		}
		length := float32(math.Max(1, math.Hypot(float64(g.joystick.X), float64(g.joystick.Y))))
		g.joystick.X, g.joystick.Y = g.joystick.X/length, g.joystick.Y/length
		g.joystick.Down = g.joystick.Down || ebiten.IsKeyPressed(ebiten.KeyShift)
	}
}

func (g *game) Draw(screen *ebiten.Image) {
	g.view.DrawBody(screen)
	g.drawInputs(screen)
	w, h := int(g.profile.DisplayWidth), int(g.profile.DisplayHeight)
	if g.oled == nil {
		g.oled = ebiten.NewImage(w, h)
	}
	for y := range h {
		for x := range w {
			i := (y*w + x) * 4
			c := byte(0)
			if g.display.Pixel(int16(x), int16(y)) {
				c = 255
			}
			g.pixels[i], g.pixels[i+1], g.pixels[i+2], g.pixels[i+3] = c, c, c, 255
		}
	}
	g.oled.WritePixels(g.pixels)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(g.view.Display.Dx())/float64(w), float64(g.view.Display.Dy())/float64(h))
	op.GeoM.Translate(float64(g.view.Display.Min.X), float64(g.view.Display.Min.Y))
	op.ColorScale.Scale(0.72, 0.84, 1, 1)
	screen.DrawImage(g.oled, op)
	x, y := float32(g.view.Encoder.X), float32(g.view.Encoder.Y)
	r := float32(g.view.EncoderRadius)
	vector.FillCircle(screen, x+2, y+7, r+5, color.RGBA{R: 13, G: 15, B: 17, A: 255}, true)
	knobColor := color.RGBA{R: 217, G: 217, B: 203, A: 255}
	if g.encoder.Pressed() || g.pointer.EncoderActive {
		knobColor = color.RGBA{R: 182, G: 185, B: 173, A: 255}
	}
	vector.FillCircle(screen, x, y, r, knobColor, true)
	for i := range 48 {
		a := float64(i) * math.Pi / 24
		cs, sn := float32(math.Cos(a)), float32(math.Sin(a))
		vector.StrokeLine(screen, x+cs*(r-5), y+sn*(r-5), x+cs*(r-1), y+sn*(r-1), 1.5,
			color.RGBA{R: 133, G: 138, B: 128, A: 255}, true)
	}
	angle := float64(g.encoder.Position()%20)*math.Pi/10 - math.Pi/2
	vector.StrokeLine(screen, x+float32(math.Cos(angle))*(r-20), y+float32(math.Sin(angle))*(r-20),
		x+float32(math.Cos(angle))*(r-9), y+float32(math.Sin(angle))*(r-9), 3,
		color.RGBA{R: 87, G: 94, B: 89, A: 255}, true)
	board.Label(screen, fmt.Sprintf("Encoder  %d", g.encoder.Position()), float64(g.view.Width-232), 47, 15, board.Ink)
	if g.encoder.Pressed() {
		board.Label(screen, "pressed", float64(g.view.Width-115), 47, 15, board.Ink)
	}
	if g.bootloader {
		vector.FillRect(screen, 178, 371, float32(g.view.Width-356), 183, color.RGBA{R: 24, G: 29, B: 33, A: 255}, false)
		board.Label(screen, "Bootloader (simulated)", 202, 395, 23, color.White)
		board.Label(screen, "Application paused. USB / UF2 flashing is not emulated.", 202, 439, 15, color.White)
		board.Label(screen, "Press RESET or F5 to start the application again.", 202, 480, 15, color.White)
	}
	if g.zoom {
		vector.FillRect(screen, 0, 0, float32(g.view.Width), float32(g.view.Height), color.RGBA{A: 215}, false)
		zx, zy := (g.view.Width-w*4)/2, (g.view.Height-h*4)/2
		vector.FillRect(screen, float32(zx-16), float32(zy-44), float32(w*4+32), float32(h*4+82), color.RGBA{R: 30, G: 32, B: 35, A: 255}, false)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(4, 4)
		op.GeoM.Translate(float64(zx), float64(zy))
		screen.DrawImage(g.oled, op)
		board.Label(screen, "OLED  /  4x", float64(zx), float64(zy-32), 16, color.White)
		board.Label(screen, "Click, O or Esc to return", float64(zx), float64(zy+h*4+12), 14, color.White)
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return g.view.Width, g.view.Height
}
