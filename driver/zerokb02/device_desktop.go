//go:build !tinygo

package zerokb02

import "github.com/rin2yh/rp2040-simulator/internal/bridge"

type hardware struct{}

var call = bridge.Call

func New() (*Device, error) { return &Device{}, nil }

func (d *Device) Read() (State, error) {
	var snapshot bridge.Inputs
	if err := call("Inputs.Read", &struct{}{}, &snapshot); err != nil {
		return State{}, err
	}
	state := State{Keys: snapshot.Keys, EncoderPressed: snapshot.EncoderPressed,
		JoystickX: snapshot.JoystickX, JoystickY: snapshot.JoystickY, JoystickPressed: snapshot.JoystickPressed}
	if d.generation == 0 || d.generation == snapshot.Generation {
		state.EncoderDelta = snapshot.EncoderPosition - d.position
	}
	d.position, d.generation = snapshot.EncoderPosition, snapshot.Generation
	return state, nil
}

func (d *Device) Display() error {
	return call("Display.WriteFrame", &bridge.DisplayFrame{Pixels: d.frame[:]}, &struct{}{})
}
