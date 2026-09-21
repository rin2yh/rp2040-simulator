//go:build !tinygo

package buzzer

import "github.com/rin2yh/rp2040-simulator/internal/bridge"

type Device struct{}

var call = bridge.Call

func newDevice() (Device, error) { return Device{}, nil }

func (Device) setFrequency(hz uint32) error {
	return call("Buzzer.SetFrequency", &bridge.BuzzerArgs{Frequency: hz}, &struct{}{})
}
