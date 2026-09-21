//go:build tinygo && xiao_rp2040

package buzzer

import "machine"

type Device struct{ channel uint8 }

func newDevice() (Device, error) {
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

func (d Device) setFrequency(hz uint32) error {
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
