package emulator

import (
	"testing"
)

func TestEncoderDetentsAndFractionalWheel(t *testing.T) {
	e := &Encoder{}
	for _, step := range []struct {
		delta float64
		want  int
	}{
		{1, 1}, {-2, -1}, {0.25, -1}, {0.25, -1}, {0.5, 0},
		{-0.75, 0}, {-0.25, -1}, {4, 3},
	} {
		e.Rotate(step.delta)
		if e.Position() != step.want {
			t.Fatalf("after %v: got %d, want %d", step.delta, e.Position(), step.want)
		}
	}
	e.SetPressed(true)
	if !e.Pressed() || e.Position() != 3 {
		t.Fatal("push must not change position")
	}
	e.SetPressed(false)
	if e.Pressed() {
		t.Fatal("release was lost")
	}
}
