//go:build !tinygo

// Package machine exposes supported board pins on desktop and TinyGo builds.
package machine

import (
	"fmt"
	"os"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/bridge"
	"github.com/rin2yh/rp2040-simulator/internal/emulator"
)

type Pin uint8
type PinMode uint8

const (
	PinOutput PinMode = iota
)

type PinConfig struct {
	Mode PinMode
}

func (Pin) Configure(PinConfig) {}

// GPIO1 drives the 12 WS2812 LEDs beneath the zero-kb02 keys.
const GPIO1 Pin = 1

// GPIO0 drives the two conf2025badge key LEDs.
const GPIO0 Pin = 0

func init() {
	handled, err := runEmulatorProcess(os.Getenv, emulator.Run)
	if !handled {
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func runEmulatorProcess(getenv func(string) string, run func(board.Profile) error) (bool, error) {
	if getenv(bridge.EmulatorProcess) != "1" {
		return false, nil
	}
	profile, err := board.Lookup(getenv(board.Environment))
	if err != nil {
		return true, err
	}
	return true, run(profile)
}
