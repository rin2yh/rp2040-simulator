//go:build !tinygo

// Package machine exposes supported board pins on desktop and TinyGo builds.
package machine

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

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

var selectedBoard = flag.String("board", "zero-kb02", "board to emulate")

func init() {
	handled, err := runEmulatorProcess(os.Getenv, os.Args[1:], emulator.Run)
	if !handled {
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func runEmulatorProcess(getenv func(string) string, args []string, run func(board.Profile) error) (bool, error) {
	if getenv(bridge.EmulatorProcess) != "1" {
		return false, nil
	}
	name := *selectedBoard
	for i, arg := range args {
		if arg == "--board" {
			if i+1 == len(args) {
				return true, errors.New("flag needs an argument: --board")
			}
			name = args[i+1]
		}
		if value, ok := strings.CutPrefix(arg, "--board="); ok {
			name = value
		}
	}
	profile, err := board.Lookup(name)
	if err != nil {
		return true, err
	}
	return true, run(profile)
}
