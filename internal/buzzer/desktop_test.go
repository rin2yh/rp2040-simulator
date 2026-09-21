//go:build !tinygo

package buzzer

import (
	"errors"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

func TestDeviceRPC(t *testing.T) {
	previous := call
	t.Cleanup(func() { call = previous })
	var frequencies []uint32
	want := errors.New("RPC unavailable")
	call = func(method string, args, reply any) error {
		if method != "Buzzer.SetFrequency" {
			t.Fatalf("method = %q", method)
		}
		frequencies = append(frequencies, args.(*bridge.BuzzerArgs).Frequency)
		return want
	}
	d, err := New(1)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.SetFrequency(440); !errors.Is(err, want) {
		t.Fatal(err)
	}
	if err := d.Stop(); !errors.Is(err, want) {
		t.Fatal(err)
	}
	for _, hz := range []uint32{1, 27, 20001, ^uint32(0)} {
		if err := d.SetFrequency(hz); err == nil {
			t.Fatalf("accepted %d Hz", hz)
		}
	}
	if len(frequencies) != 2 || frequencies[0] != 440 || frequencies[1] != 0 {
		t.Fatal(frequencies)
	}
	if _, err := New(0); err == nil {
		t.Fatal("accepted unsupported pin")
	}
}
