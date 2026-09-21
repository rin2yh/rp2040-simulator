//go:build !tinygo

// Package ws2812 provides the same application-facing operations on desktop
// and supported TinyGo hardware.
package ws2812

import (
	"errors"
	"image/color"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
	"github.com/rin2yh/rp2040-simulator/machine"
)

type Device struct {
	pin machine.Pin
}

var call = bridge.Call

func NewWS2812(pin machine.Pin) Device { return Device{pin: pin} }

func (d Device) WriteColors(colors []color.RGBA) error {
	if d.pin != machine.GPIO1 && d.pin != machine.GPIO0 {
		return errors.New("unsupported WS2812 pin")
	}
	args := bridge.WriteLEDsArgs{Colors: append([]color.RGBA(nil), colors...)}
	return call("LEDs.WriteColors", &args, &struct{}{})
}
