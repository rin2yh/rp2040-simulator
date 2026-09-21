// Package buzzer implements the shared logical buzzer operations.
package buzzer

import "errors"

// ValidateFrequency accepts silence or the supported audible frequency range.
func ValidateFrequency(hz uint32) error {
	if hz != 0 && (hz < 28 || hz > 20000) {
		return errors.New("buzzer frequency must be 0 (silence) or 28..20000 Hz")
	}
	return nil
}
