//go:build tinygo && waveshare_rp2040_zero

package zerokb02

import (
	"image/color"
	"machine"
	"runtime/interrupt"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/ssd1306"
)

type hardware struct {
	columns                           [4]machine.Pin
	rows                              [3]machine.Pin
	encoderA, encoderB, encoderButton machine.Pin
	joystickButton                    machine.Pin
	x, y                              machine.ADC
	previousEncoder                   uint8
	encoderSteps                      int
	display                           *ssd1306.Device
}

// New configures the board's matrix, encoder, ADC, and SSD1306 display.
func New() (*Device, error) {
	d := &Device{}
	h := &d.hardware
	h.columns = [4]machine.Pin{machine.GPIO5, machine.GPIO6, machine.GPIO7, machine.GPIO8}
	h.rows = [3]machine.Pin{machine.GPIO9, machine.GPIO10, machine.GPIO11}
	for _, p := range h.columns {
		p.Configure(machine.PinConfig{Mode: machine.PinOutput})
		p.High()
	}
	for _, p := range h.rows {
		p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	h.encoderA, h.encoderB, h.encoderButton = machine.GPIO3, machine.GPIO4, machine.GPIO2
	h.joystickButton = machine.GPIO0
	for _, p := range [...]machine.Pin{h.encoderA, h.encoderB, h.encoderButton, h.joystickButton} {
		p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}
	h.previousEncoder = h.encoderBits()
	machine.InitADC()
	h.x, h.y = machine.ADC{Pin: machine.GPIO29}, machine.ADC{Pin: machine.GPIO28}
	if err := h.x.Configure(machine.ADCConfig{}); err != nil {
		return nil, err
	}
	if err := h.y.Configure(machine.ADCConfig{}); err != nil {
		return nil, err
	}
	if err := machine.I2C0.Configure(machine.I2CConfig{SDA: machine.GPIO12, SCL: machine.GPIO13, Frequency: 400000}); err != nil {
		return nil, err
	}
	h.display = ssd1306.NewI2C(machine.I2C0)
	h.display.Configure(ssd1306.Config{Width: Width, Height: Height, Address: 0x3c, Rotation: drivers.Rotation180})
	if err := h.encoderA.SetInterrupt(machine.PinToggle, h.encoderEdge); err != nil {
		return nil, err
	}
	if err := h.encoderB.SetInterrupt(machine.PinToggle, h.encoderEdge); err != nil {
		_ = h.encoderA.SetInterrupt(0, nil)
		return nil, err
	}
	return d, nil
}

func (h *hardware) encoderBits() uint8 {
	var bits uint8
	if h.encoderA.Get() {
		bits |= 1
	}
	if h.encoderB.Get() {
		bits |= 2
	}
	return bits
}

func (h *hardware) encoderEdge(machine.Pin) {
	current := h.encoderBits()
	transitions := [16]int8{0, -1, 1, 0, 1, 0, 0, -1, -1, 0, 0, 1, 0, 1, -1, 0}
	h.encoderSteps += int(transitions[h.previousEncoder<<2|current])
	h.previousEncoder = current
}

func axis(raw uint16) float32 {
	const low, center, high = 0x3000, 0x8000, 0xc800
	if raw < center {
		if raw <= low {
			return -1
		}
		return float32(int(raw)-center) / float32(center-low)
	}
	if raw >= high {
		return 1
	}
	return float32(int(raw)-center) / float32(high-center)
}

func (d *Device) Read() (State, error) {
	h := &d.hardware
	var state State
	for col, p := range h.columns {
		p.Low()
		time.Sleep(time.Millisecond)
		for row, input := range h.rows {
			state.Keys[row*4+col] = !input.Get()
		}
		p.High()
	}
	// Four edges form one detent. Preserve a fractional turn across reads.
	mask := interrupt.Disable()
	state.EncoderDelta = h.encoderSteps / 4
	h.encoderSteps -= state.EncoderDelta * 4
	interrupt.Restore(mask)
	state.EncoderPressed = !h.encoderButton.Get()
	state.JoystickX, state.JoystickY = axis(h.x.Get()), -axis(h.y.Get())
	state.JoystickPressed = !h.joystickButton.Get()
	return state, nil
}

func (d *Device) Display() error {
	h := d.hardware.display
	h.ClearBuffer()
	for y := int16(0); y < Height; y++ {
		for x := int16(0); x < Width; x++ {
			if d.frame[int(x)+int(y/8)*Width]&(1<<uint(y%8)) != 0 {
				h.SetPixel(x, y, color.RGBA{R: 255})
			}
		}
	}
	return h.Display()
}
