//go:build tinygo && (waveshare_rp2040_zero || xiao_rp2040)

// Package machine exposes supported board pins on desktop and TinyGo builds.
package machine

import m "machine"

type Pin = m.Pin
type PinMode = m.PinMode
type PinConfig = m.PinConfig

const PinOutput = m.PinOutput

const GPIO1 = m.GPIO1
const GPIO0 = m.GPIO0
