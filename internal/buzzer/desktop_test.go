//go:build !tinygo

package buzzer

import (
	"errors"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
	"github.com/rin2yh/rp2040-simulator/internal/tabletest"
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
	if len(frequencies) != 2 || frequencies[0] != 440 || frequencies[1] != 0 {
		t.Fatal(frequencies)
	}
	if _, err := New(0); err == nil {
		t.Fatal("accepted unsupported pin")
	}
}

func TestValidateFrequency(t *testing.T) {
	type want struct {
		err bool
	}
	cases := []tabletest.Case[uint32, want]{
		{Name: "silence", In: 0},
		{Name: "below minimum", In: 27, Want: want{err: true}},
		{Name: "minimum", In: 28},
		{Name: "maximum", In: 20000},
		{Name: "above maximum", In: 20001, Want: want{err: true}},
		{Name: "uint32 maximum", In: ^uint32(0), Want: want{err: true}},
	}

	tabletest.Run(t, cases, func(t *testing.T, hz uint32, want want) {
		err := ValidateFrequency(hz)
		if (err != nil) != want.err {
			t.Fatalf("ValidateFrequency(%d) error = %v, want error %v", hz, err, want.err)
		}
	})
}
