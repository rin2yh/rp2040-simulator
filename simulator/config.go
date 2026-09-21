// Package simulator configures the desktop simulator used by an application.
package simulator

import "fmt"

// Config controls the simulator process that is started on the first RPC call.
type Config struct {
	Board string
}

// Configure sets simulator options. It must be called before using a simulated
// device. An empty board selects zero-kb02.
func Configure(config Config) error {
	name := config.Board
	if name == "" {
		name = "zero-kb02"
	}
	if name != "zero-kb02" && name != "conf2025badge" {
		return fmt.Errorf("unknown board %q (want zero-kb02 or conf2025badge)", name)
	}
	return configureBoard(name)
}
