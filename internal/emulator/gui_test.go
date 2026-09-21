package emulator

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rin2yh/rp2040-simulator/internal/board"
)

var guiError error

var errGUICheckFinished = errors.New("GUI check finished")

type guiCheck struct {
	profiles []board.Profile
	rendered bool
	err      error
}

func (c *guiCheck) Update() error {
	if c.rendered {
		return errGUICheckFinished
	}
	return nil
}

func (c *guiCheck) Draw(screen *ebiten.Image) {
	c.rendered = true
	for _, profile := range c.profiles {
		actual, err := renderBoard(profile)
		if err != nil {
			c.err = fmt.Errorf("render %s: %w", profile.Name, err)
			return
		}
		if err := checkBoardRendering(actual, profile); err != nil {
			c.err = fmt.Errorf("check %s: %w", profile.Name, err)
			return
		}
		screen.DrawImage(actual, nil)
	}
}

func (c *guiCheck) Layout(_, _ int) (int, int) { return 890, 860 }

// Ebitengine must run on the main goroutine, hence TestMain. Enable explicitly
// with the gui build tag on a desktop or under Xvfb. Ordinary tests remain headless.
func TestMain(m *testing.M) {
	guiError = checkGUI()
	os.Exit(m.Run())
}

func TestGUIRenderingMatchesBoardGolden(t *testing.T) {
	if guiError != nil {
		t.Fatal(guiError)
	}
}

func checkGUI() error {
	check := &guiCheck{profiles: []board.Profile{
		board.ZeroKB02(),
		board.Conf2025Badge(),
	}}
	ebiten.SetWindowSize(890, 860)
	ebiten.SetWindowTitle("TinyGo device emulator GUI check")
	err := ebiten.RunGame(check)
	if !errors.Is(err, errGUICheckFinished) {
		return fmt.Errorf("GUI exited before check completed: %v", err)
	}
	return check.err
}

func renderBoard(profile board.Profile) (*ebiten.Image, error) {
	g, err := newGame(profile)
	if err != nil {
		return nil, err
	}
	g.hint = defaultHint
	g.encoder.Rotate(-2)
	g.leds.Set(0, color.RGBA{R: 255})
	g.leds.Set(10, color.RGBA{B: 255})
	g.leds.Set(min(5, g.leds.Len()-1), color.RGBA{G: 255})
	if err := g.leds.Display(); err != nil {
		return nil, err
	}
	for x := int16(0); x < profile.DisplayWidth; x++ {
		g.display.SetPixel(x, x%profile.DisplayHeight, color.RGBA{R: 255})
	}
	if err := g.display.Display(); err != nil {
		return nil, err
	}
	g.buttons[0].Down = true
	g.joystick.X, g.joystick.Y = 0.65, -0.4
	actual := ebiten.NewImage(profile.View.Width, profile.View.Height)
	g.Draw(actual)
	return actual, nil
}

func checkBoardRendering(actual *ebiten.Image, profile board.Profile) error {
	view := profile.View
	lit := 0
	for y := view.Display.Min.Y; y < view.Display.Max.Y; y++ {
		for x := view.Display.Min.X; x < view.Display.Max.X; x++ {
			r, _, _, _ := actual.At(x, y).RGBA()
			if r > 0x8000 {
				lit++
			}
		}
	}
	if lit == 0 || lit == view.Display.Dx()*view.Display.Dy() {
		return fmt.Errorf("OLED rendering is blank: %d lit pixels", lit)
	}
	key := view.Keys[0].Bounds
	_, green, blue, _ := actual.At(key.Min.X+3, key.Min.Y+50).RGBA()
	if green < 0xc000 || blue < 0xc000 {
		return fmt.Errorf("pressed key highlight missing: green=%x blue=%x", green, blue)
	}
	led := view.LEDs[min(5, len(view.LEDs)-1)]
	red, green, blue, _ := actual.At(led.X, led.Y).RGBA()
	if green < 0xc000 || green <= red || green <= blue {
		return fmt.Errorf("RGB LED rendering missing: red=%x green=%x blue=%x", red, green, blue)
	}
	return checkBoardGolden(actual, profile.Name)
}

func checkBoardGolden(screen *ebiten.Image, name string) error {
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, screen); err != nil {
		return fmt.Errorf("encode screenshot: %w", err)
	}
	actualPNG := encoded.Bytes()
	if path := os.Getenv("EMULATOR_SNAPSHOT"); path != "" {
		if err := os.WriteFile(path, actualPNG, 0o644); err != nil {
			return fmt.Errorf("write screenshot: %w", err)
		}
	}
	golden := boardGoldenPath(name)
	if os.Getenv("UPDATE_BOARD_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(golden), 0o755); err != nil {
			return fmt.Errorf("create board testdata: %w", err)
		}
		if err := os.WriteFile(golden, actualPNG, 0o644); err != nil {
			return fmt.Errorf("update board golden: %w", err)
		}
		return nil
	}
	expectedFile, err := os.Open(golden)
	if err != nil {
		return fmt.Errorf("open board golden %q: %w", golden, err)
	}
	defer expectedFile.Close()
	expected, err := png.Decode(expectedFile)
	if err != nil {
		return fmt.Errorf("decode board golden: %w", err)
	}
	actual, err := png.Decode(bytes.NewReader(actualPNG))
	if err != nil {
		return fmt.Errorf("decode screenshot: %w", err)
	}
	return compareImages(actual, expected)
}

func boardGoldenPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "board", "testdata", name+".png")
}

func compareImages(actual, expected image.Image) error {
	if actual.Bounds() != expected.Bounds() {
		return fmt.Errorf("board golden size = %v, screenshot size = %v", expected.Bounds(), actual.Bounds())
	}
	different := 0
	for y := actual.Bounds().Min.Y; y < actual.Bounds().Max.Y; y++ {
		for x := actual.Bounds().Min.X; x < actual.Bounds().Max.X; x++ {
			if color.NRGBAModel.Convert(actual.At(x, y)) != color.NRGBAModel.Convert(expected.At(x, y)) {
				different++
			}
		}
	}
	if different != 0 {
		return fmt.Errorf("board rendering differs from golden at %d pixels; inspect build/emulator.png", different)
	}
	return nil
}
