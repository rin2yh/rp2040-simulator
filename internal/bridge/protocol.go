// Package bridge defines the private protocol between desktop drivers and the
// emulator process.
package bridge

import "image/color"

const (
	Address         = "127.0.0.1:9840"
	EmulatorProcess = "RP2040_SIMULATOR_EMULATOR_PROCESS"
)

type WriteLEDsArgs struct {
	Colors []color.RGBA
}

type BuzzerArgs struct {
	Frequency uint32
}
