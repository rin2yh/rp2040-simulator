//go:build tinygo && xiao_rp2040

package buzzer

import (
	"errors"
	"machine"
)

type Device struct{ channel uint8 }

func New(pin uint8) (Device, error) {
	if machine.Pin(pin) != machine.GPIO1 {
		return Device{}, errors.New("conf2025badge buzzer requires GPIO1")
	}
	if err := machine.PWM0.Configure(machine.PWMConfig{Period: 1000000000 / 440}); err != nil {
		return Device{}, err
	}
	ch, err := machine.PWM0.Channel(machine.GPIO1)
	if err != nil {
		return Device{}, err
	}
	machine.PWM0.Set(ch, 0)
	return Device{channel: ch}, nil
}

func (d Device) SetFrequency(hz uint32) error {
	if err := ValidateFrequency(hz); err != nil {
		return err
	}
	machine.PWM0.Set(d.channel, 0)
	if hz == 0 {
		return nil
	}
	if err := machine.PWM0.SetPeriod(1000000000 / uint64(hz)); err != nil {
		return err
	}
	machine.PWM0.Set(d.channel, machine.PWM0.Top()/2)
	return nil
}

func (d Device) Stop() error { return d.SetFrequency(0) }
