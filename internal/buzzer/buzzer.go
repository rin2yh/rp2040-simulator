// Package buzzer implements the shared logical buzzer operations.
package buzzer

import "errors"

func New(pin uint8) (Device, error) {
	if pin != 1 {
		return Device{}, errors.New("conf2025badge buzzer requires GPIO1")
	}
	return newDevice()
}

func (d Device) SetFrequency(hz uint32) error {
	if err := ValidateFrequency(hz); err != nil {
		return err
	}
	return d.setFrequency(hz)
}

func (d Device) Stop() error { return d.SetFrequency(0) }

// ValidateFrequency accepts silence or the supported audible frequency range.
func ValidateFrequency(hz uint32) error {
	if hz != 0 && (hz < 28 || hz > 20000) {
		return errors.New("buzzer frequency must be 0 (silence) or 28..20000 Hz")
	}
	return nil
}
