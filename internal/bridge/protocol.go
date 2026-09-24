// Package bridge defines the private protocol between desktop drivers and the
// emulator process.
package bridge

import "image/color"

const (
	Address         = "127.0.0.1:9840"
	EmulatorProcess = "RP2040_SIMULATOR_EMULATOR_PROCESS"
	EmulatorBoard   = "RP2040_SIMULATOR_EMULATOR_BOARD"
)

type WriteLEDsArgs struct {
	Colors []color.RGBA
}

type BuzzerArgs struct {
	Frequency uint32
}

// Inputs is a snapshot of zero-kb02 controls. Position is cumulative so a
// client can calculate detents without consuming another client's events.
type Inputs struct {
	Keys                 [12]bool
	EncoderPosition      int
	EncoderPressed       bool
	JoystickX, JoystickY float32
	JoystickPressed      bool
	Generation           uint64
}

type DisplayFrame struct {
	// Pixels are packed in SSD1306 page order (128 columns by 8 pages).
	Pixels []byte
}
