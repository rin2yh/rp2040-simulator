// Package simulator configures the desktop simulator used by an application.
package simulator

import "fmt"

// Board identifies a board supported by the simulator.
type Board string

const (
	BoardZeroKB02      Board = "zero-kb02"
	BoardConf2025Badge Board = "conf2025badge"
)

// Config controls the simulator process that is started on the first RPC call.
type Config struct {
	Board Board
}

// Configure sets simulator options. It must be called before using a simulated
// device.
func Configure(config Config) error {
	switch config.Board {
	case BoardZeroKB02, BoardConf2025Badge:
		return configureBoard(string(config.Board))
	default:
		return fmt.Errorf("unknown board %q (want %q or %q)", config.Board, BoardZeroKB02, BoardConf2025Badge)
	}
}
