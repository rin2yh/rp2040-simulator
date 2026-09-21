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
	tests := []struct {
		name string
		hz   uint32
		want want
	}{
		{name: "silence", hz: 0},
		{name: "below minimum", hz: 27, want: want{err: true}},
		{name: "minimum", hz: 28},
		{name: "maximum", hz: 20000},
		{name: "above maximum", hz: 20001, want: want{err: true}},
		{name: "uint32 maximum", hz: ^uint32(0), want: want{err: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFrequency(tt.hz)
			if (err != nil) != tt.want.err {
				t.Fatalf("ValidateFrequency(%d) error = %v, want error %v", tt.hz, err, tt.want.err)
			}
		})
	}
}
