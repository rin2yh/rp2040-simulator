//go:build tinygo && waveshare_rp2040_zero

package main

import (
	"image/color"

	"github.com/rin2yh/rp2040-simulator/driver/ws2812"
	"github.com/rin2yh/rp2040-simulator/machine"
)

func main() {
	machine.GPIO1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.NewWS2812(machine.GPIO1)
	_ = leds.WriteColors([]color.RGBA{{}})
}
