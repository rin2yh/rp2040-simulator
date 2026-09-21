//go:build !tinygo

package machine

import (
	"errors"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

func TestRunEmulatorProcessIgnoresApplicationProcess(t *testing.T) {
	called := false
	handled, err := runEmulatorProcess(func(string) string { return "" }, func(board.Profile) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if handled || called {
		t.Fatal("application process started the emulator entry point")
	}
}

func TestRunEmulatorProcessSelectsBoard(t *testing.T) {
	type want struct {
		profile string
		leds    int
		keys    int
		called  bool
		err     bool
	}
	tests := []struct {
		name  string
		board string
		want  want
	}{
		{name: "default", want: want{profile: "zero-kb02", leds: 12, keys: 12, called: true, err: true}},
		{name: "conf2025badge", board: "conf2025badge", want: want{profile: "conf2025badge", leds: 2, keys: 2, called: true, err: true}},
		{name: "unknown", board: "unknown", want: want{err: true}},
	}

	runErr := errors.New("stopped")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			handled, err := runEmulatorProcess(func(key string) string {
				if key == bridge.EmulatorProcess {
					return "1"
				}
				if key == bridge.EmulatorBoard {
					return tt.board
				}
				return ""
			}, func(p board.Profile) error {
				called = true
				if p.Name != tt.want.profile || p.LEDCount != tt.want.leds || len(p.View.Keys) != tt.want.keys {
					t.Fatalf("profile = {name: %q, LEDs: %d, keys: %d}, want %+v", p.Name, p.LEDCount, len(p.View.Keys), tt.want)
				}
				return runErr
			})
			if !handled {
				t.Fatal("emulator process was not handled")
			}
			if (err != nil) != tt.want.err {
				t.Fatalf("error = %v, want error %v", err, tt.want.err)
			}
			if called != tt.want.called {
				t.Fatalf("called = %v, want %v", called, tt.want.called)
			}
			if called && !errors.Is(err, runErr) {
				t.Fatalf("error = %v, want %v", err, runErr)
			}
		})
	}
}
