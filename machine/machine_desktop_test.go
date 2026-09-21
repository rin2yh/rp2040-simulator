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

func TestRunEmulatorProcessStartsZeroKB02(t *testing.T) {
	want := errors.New("stopped")
	handled, err := runEmulatorProcess(func(name string) string {
		if name == bridge.EmulatorProcess {
			return "1"
		}
		return ""
	}, func(profile board.Profile) error {
		if profile.Name != "zero-kb02" {
			t.Fatalf("profile = %q, want zero-kb02", profile.Name)
		}
		return want
	})
	if !handled {
		t.Fatal("emulator process was not handled")
	}
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want %v", err, want)
	}
}
