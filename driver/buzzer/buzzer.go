//go:build !tinygo || xiao_rp2040

// Package buzzer controls the conf2025badge speaker on desktop and TinyGo.
package buzzer

import (
	"github.com/rin2yh/rp2040-simulator/internal/buzzer"
	"github.com/rin2yh/rp2040-simulator/machine"
)

// Device plays a continuous square wave until Stop or another SetFrequency.
// SetFrequency accepts 28..20000 Hz, or zero for silence.
type Device = buzzer.Device

// New configures the conf2025badge buzzer on machine.GPIO1.
func New(pin machine.Pin) (Device, error) { return buzzer.New(uint8(pin)) }
