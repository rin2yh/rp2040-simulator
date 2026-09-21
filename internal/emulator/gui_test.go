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
	profile := board.ZeroKB02()
	view := profile.View
	drawBody := view.DrawBody
	draws, lastStep := 0, -1
	var captureError error
	finished := errors.New("GUI check finished")
	// Preserve the previous completed frame so the next DrawBody can inspect
	// the actual composed screen before it is painted over.
	ebiten.SetScreenClearedEveryFrame(false)
	defer ebiten.SetScreenClearedEveryFrame(true)
	view.DrawBody = func(screen *ebiten.Image) {
		if draws == 5 {
			lit := 0
			for y := view.Display.Min.Y; y < view.Display.Max.Y; y++ {
				for x := view.Display.Min.X; x < view.Display.Max.X; x++ {
					r, _, _, _ := screen.At(x, y).RGBA()
					if r > 0x8000 {
						lit++
					}
				}
			}
			if lit == 0 || lit == view.Display.Dx()*view.Display.Dy() {
				captureError = fmt.Errorf("OLED rendering is blank: %d lit pixels", lit)
			}
			key := view.Keys[0].Bounds
			_, green, blue, _ := screen.At(key.Min.X+3, key.Min.Y+50).RGBA()
			if green < 0xc000 || blue < 0xc000 {
				captureError = fmt.Errorf("pressed key highlight missing: green=%x blue=%x", green, blue)
			}
			led := view.LEDs[5]
			red, green, blue, _ := screen.At(led.X, led.Y).RGBA()
			if green < 0xc000 || green <= red || green <= blue {
				captureError = fmt.Errorf("RGB LED rendering missing: red=%x green=%x blue=%x", red, green, blue)
			}
			if err := checkBoardGolden(screen, profile.Name); err != nil {
				captureError = err
			}
		}
		draws++
		drawBody(screen)
	}
	profile.View = view
	g, err := newGame(profile)
	if err != nil {
		return err
	}
	g.testUpdate = func(g *game) error {
		// Keep the golden image independent of the host pointer position.
		g.hint = defaultHint
		if draws > 5 {
			return finished
		}
		if draws != lastStep {
			switch draws {
			case 1:
				g.encoder.Rotate(3)
				g.leds.Set(0, color.RGBA{R: 255})
				g.leds.Set(5, color.RGBA{G: 255})
				g.leds.Set(10, color.RGBA{B: 255})
				if err := g.leds.Display(); err != nil {
					return err
				}
				for x := int16(0); x < 128; x++ {
					g.display.SetPixel(x, x%64, color.RGBA{R: 255})
				}
				if err := g.display.Display(); err != nil {
					return err
				}
			case 2:
				g.encoder.Rotate(-5)
			case 3:
				g.encoder.SetPressed(true)
			case 4:
				g.encoder.SetPressed(false)
			}
			lastStep = draws
		}
		// Hold through every Update between these Draw calls, just as the
		// input sampler does for real keyboard and pointer gestures.
		if draws == 4 {
			g.buttons[0].Down = true
			g.joystick.X, g.joystick.Y = 0.65, -0.4
		}
		return nil
	}
	err = run(g)
	if !errors.Is(err, finished) {
		return fmt.Errorf("GUI exited before check completed: %v", err)
	}
	return captureError
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
