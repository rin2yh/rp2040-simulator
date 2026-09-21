//go:build !tinygo

package buzzer

import (
	"errors"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

type Device struct{}

var call = bridge.Call

func New(pin uint8) (Device, error) {
	if pin != 1 {
		return Device{}, errors.New("conf2025badge buzzer requires GPIO1")
	}
	return Device{}, nil
}

func (Device) SetFrequency(hz uint32) error {
	if err := ValidateFrequency(hz); err != nil {
		return err
	}
	return call("Buzzer.SetFrequency", &bridge.BuzzerArgs{Frequency: hz}, &struct{}{})
}

func (d Device) Stop() error { return d.SetFrequency(0) }
