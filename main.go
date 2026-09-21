package main

import (
	"log"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/emulator"
)

func main() {
	if err := emulator.Run(board.ZeroKB02()); err != nil {
		log.Fatal(err)
	}
}
