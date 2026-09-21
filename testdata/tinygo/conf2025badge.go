//go:build tinygo && xiao_rp2040

package main

import (
	"image/color"

	"github.com/rin2yh/rp2040-simulator/driver/buzzer"
	"github.com/rin2yh/rp2040-simulator/driver/ws2812"
	"github.com/rin2yh/rp2040-simulator/machine"
	"github.com/rin2yh/rp2040-simulator/simulator"
)

func main() {
	if err := simulator.Configure(simulator.Config{Board: "conf2025badge"}); err != nil {
		panic(err)
	}
	speaker, err := buzzer.New(machine.GPIO1)
	if err != nil {
		panic(err)
	}
	if err := speaker.SetFrequency(440); err != nil {
		panic(err)
	}
	if err := speaker.Stop(); err != nil {
		panic(err)
	}
	machine.GPIO0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.NewWS2812(machine.GPIO0)
	_ = leds.WriteColors([]color.RGBA{{R: 255}, {B: 255}})
}
