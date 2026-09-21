//go:build tinygo && waveshare_rp2040_zero

package ws2812

import (
	"github.com/rin2yh/rp2040-simulator/machine"
	tinyws2812 "tinygo.org/x/drivers/ws2812"
)

type Device = tinyws2812.Device

func NewWS2812(pin machine.Pin) Device { return tinyws2812.NewWS2812(pin) }
