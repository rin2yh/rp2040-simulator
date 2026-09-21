package main

import (
	"image/color"
	"log"
	"time"

	"github.com/rin2yh/rp2040-simulator/driver/ws2812"
	"github.com/rin2yh/rp2040-simulator/machine"
)

const ledCount = 12

func main() {
	machine.GPIO1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	leds := ws2812.NewWS2812(machine.GPIO1)
	colors := make([]color.RGBA, ledCount)

	for {
		fill(colors, color.RGBA{G: 24})
		if err := leds.WriteColors(colors); err != nil {
			log.Fatal(err)
		}
		time.Sleep(500 * time.Millisecond)

		fill(colors, color.RGBA{})
		if err := leds.WriteColors(colors); err != nil {
			log.Fatal(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func fill(leds []color.RGBA, c color.RGBA) {
	for i := range leds {
		leds[i] = c
	}
}
