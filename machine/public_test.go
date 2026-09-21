//go:build !tinygo

package machine_test

import (
	"testing"

	"github.com/rin2yh/rp2040-simulator/machine"
)

func TestZeroKB02GPIO1(t *testing.T) {
	if machine.GPIO1 != machine.Pin(1) {
		t.Fatalf("GPIO1 = %d, want 1", machine.GPIO1)
	}
	machine.GPIO1.Configure(machine.PinConfig{Mode: machine.PinOutput})
}
