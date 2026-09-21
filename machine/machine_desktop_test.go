//go:build !tinygo

package machine

import (
	"errors"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/bridge"
	"github.com/rin2yh/rp2040-simulator/internal/tabletest"
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
	cases := []tabletest.Case[string, want]{
		{Name: "default", Want: want{profile: "zero-kb02", leds: 12, keys: 12, called: true, err: true}},
		{Name: "conf2025badge", In: "conf2025badge", Want: want{profile: "conf2025badge", leds: 2, keys: 2, called: true, err: true}},
		{Name: "unknown", In: "unknown", Want: want{err: true}},
	}

	runErr := errors.New("stopped")
	tabletest.Run(t, cases, func(t *testing.T, boardName string, want want) {
		called := false
		handled, err := runEmulatorProcess(func(key string) string {
			if key == bridge.EmulatorProcess {
				return "1"
			}
			if key == bridge.EmulatorBoard {
				return boardName
			}
			return ""
		}, func(p board.Profile) error {
			called = true
			if p.Name != want.profile || p.LEDCount != want.leds || len(p.View.Keys) != want.keys {
				t.Fatalf("profile = {name: %q, LEDs: %d, keys: %d}, want %+v", p.Name, p.LEDCount, len(p.View.Keys), want)
			}
			return runErr
		})
		if !handled {
			t.Fatal("emulator process was not handled")
		}
		if (err != nil) != want.err {
			t.Fatalf("error = %v, want error %v", err, want.err)
		}
		if called != want.called {
			t.Fatalf("called = %v, want %v", called, want.called)
		}
		if called && !errors.Is(err, runErr) {
			t.Fatalf("error = %v, want %v", err, runErr)
		}
	})
}
