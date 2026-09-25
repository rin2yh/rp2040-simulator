//go:build tinygo && waveshare_rp2040_zero

package main

import (
	"image/color"

	"github.com/rin2yh/rp2040-simulator/driver/ws2812"
	"github.com/rin2yh/rp2040-simulator/driver/zerokb02"
	"github.com/rin2yh/rp2040-simulator/machine"
)

func main() {
	machine.GPIO1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.NewWS2812(machine.GPIO1)
	_ = leds.WriteColors([]color.RGBA{{}})
	device, err := zerokb02.New()
	if err != nil {
		panic(err)
	}
	if _, err := device.Read(); err != nil {
		panic(err)
	}
	device.SetPixel(0, 0, color.RGBA{R: 255})
	if err := device.Display(); err != nil {
		panic(err)
	}
}
