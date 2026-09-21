//go:build tinygo && waveshare_rp2040_zero

// Package machine exposes zero-kb02 pins on desktop and TinyGo builds.
package machine

import m "machine"

type Pin = m.Pin
type PinMode = m.PinMode
type PinConfig = m.PinConfig

const PinOutput = m.PinOutput

const GPIO1 = m.GPIO1
