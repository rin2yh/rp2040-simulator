//go:build !tinygo

package simulator

import "github.com/rin2yh/rp2040-simulator/internal/bridge"

func configureBoard(name string) error {
	return bridge.ConfigureBoard(name)
}
