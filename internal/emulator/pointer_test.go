package emulator

import (
	"image"
	"math"
	"testing"
)

func testPointer() *Pointer {
	return NewPointer(PointerLayout{Keys: []image.Rectangle{image.Rect(0, 0, 40, 40), image.Rect(50, 0, 90, 40)},
		Encoder: image.Pt(150, 50), EncoderRadius: 40, Joystick: image.Pt(250, 50), JoystickRadius: 40})
}

func step(p *Pointer, x, y int, down bool) {
	p.Step(PointerInput{Position: image.Pt(x, y), Down: down, Enabled: true})
}

func TestKeyCaptureReleaseAndFocusLoss(t *testing.T) {
	p := testPointer()
	step(p, 20, 20, true)
	if !p.Keys[0] {
		t.Fatal("key does not respond to mouse down")
	}
	step(p, 60, 20, true)
	if p.Keys[0] || p.Keys[1] {
		t.Fatal("dragging switched key capture")
	}
	step(p, 20, 20, false)
	if p.Keys[0] {
		t.Fatal("key stuck after release")
	}
	step(p, 20, 20, true)
	p.Step(PointerInput{Position: image.Pt(20, 20), Down: true, Enabled: false})
	step(p, 20, 20, true)
	if p.Keys[0] {
		t.Fatal("focus loss must cancel until release")
	}
}

func TestEncoderClickDragAndLocalWheel(t *testing.T) {
	p := testPointer()
	step(p, 150, 50, true)
	step(p, 150, 50, false)
	if !p.EncoderDown || p.Delta != 0 {
		t.Fatal("click must push without rotating")
	}
	step(p, 150, 50, false)
	if p.EncoderDown {
		t.Fatal("click pulse stuck")
	}
	step(p, 150, 50, true)
	step(p, 150, 26, true)
	if p.Delta != 3 {
		t.Fatalf("center drag got %v, want 3", p.Delta)
	}
	step(p, 150, 26, false)
	if p.EncoderDown {
		t.Fatal("rotation must not produce a click/reset")
	}
	p.Step(PointerInput{Position: image.Pt(20, 20), Wheel: 3, Enabled: true})
	if p.Delta != 0 {
		t.Fatal("scroll outside knob rotated it")
	}
	p.Step(PointerInput{Position: image.Pt(150, 50), Wheel: 0.25, Enabled: true})
	if p.Delta != 0.25 {
		t.Fatal("fractional wheel lost")
	}
	// Crossing the -pi/+pi boundary is a short clockwise rotation.
	step(p, 111, 51, true)
	step(p, 111, 45, true)
	if p.Delta <= 0 || p.Delta > 1 {
		t.Fatalf("angle wrapping jumped: %v", p.Delta)
	}
}

func TestJoystickClampsSpringsAndDoesNotClickAfterDrag(t *testing.T) {
	p := testPointer()
	step(p, 250, 50, true)
	step(p, 330, 130, true)
	if math.Abs(math.Hypot(float64(p.X), float64(p.Y))-1) > 0.0001 {
		t.Fatal("joystick not clamped to unit circle")
	}
	step(p, 330, 130, false)
	if p.X != 0 || p.Y != 0 || p.JoystickDown {
		t.Fatal("release must center without pushing")
	}
	step(p, 250, 50, true)
	step(p, 250, 50, false)
	if !p.JoystickDown {
		t.Fatal("tap must push")
	}
	// Opening a modal during a gesture must not leave an axis held.
	step(p, 250, 50, true)
	step(p, 280, 50, true)
	p.Step(PointerInput{Down: true})
	if p.X != 0 || p.Y != 0 {
		t.Fatal("blocked pointer left joystick held")
	}
}
