package main

import (
	"time"

	"github.com/rin2yh/rp2040-simulator/driver/buzzer"
	"github.com/rin2yh/rp2040-simulator/machine"
	"github.com/rin2yh/rp2040-simulator/simulator"
)

func main() {
	if err := simulator.Configure(simulator.Config{Board: simulator.BoardConf2025Badge}); err != nil {
		panic(err)
	}
	speaker, err := buzzer.New(machine.GPIO1)
	if err != nil {
		panic(err)
	}
	for _, hz := range []uint32{440, 554, 659} {
		if err := speaker.SetFrequency(hz); err != nil {
			panic(err)
		}
		time.Sleep(250 * time.Millisecond)
	}
	if err := speaker.Stop(); err != nil {
		panic(err)
	}
}
