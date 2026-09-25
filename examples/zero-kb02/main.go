package main

import (
	"image/color"
	"time"

	"github.com/rin2yh/rp2040-simulator/driver/ws2812"
	"github.com/rin2yh/rp2040-simulator/driver/zerokb02"
	"github.com/rin2yh/rp2040-simulator/machine"
)

func main() {
	device, err := zerokb02.New()
	if err != nil {
		panic(err)
	}
	machine.GPIO1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.NewWS2812(machine.GPIO1)
	var colors [12]color.RGBA
	position := 64
	for {
		state, err := device.Read()
		if err != nil {
			panic(err)
		}
		position = (position + state.EncoderDelta + 128) % 128
		device.ClearBuffer()
		for i, pressed := range state.Keys {
			if pressed {
				colors[i] = color.RGBA{G: 80}
				for x := int16(i * 10); x < int16(i*10+8); x++ {
					device.SetPixel(x, 8, color.RGBA{R: 255})
				}
			} else {
				colors[i] = color.RGBA{}
			}
		}
		device.SetPixel(int16(position), 30, color.RGBA{R: 255})
		device.SetPixel(int16(64+state.JoystickX*50), int16(44+state.JoystickY*16), color.RGBA{R: 255})
		if state.EncoderPressed {
			device.SetPixel(0, 63, color.RGBA{R: 255})
		}
		if state.JoystickPressed {
			device.SetPixel(127, 63, color.RGBA{R: 255})
		}
		if err := device.Display(); err != nil {
			panic(err)
		}
		if err := leds.WriteColors(colors[:]); err != nil {
			panic(err)
		}
		time.Sleep(20 * time.Millisecond)
	}
}
