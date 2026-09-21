package main

import (
	"log"
	"os"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/emulator"
)

func main() {
	profile, err := board.Lookup(os.Getenv(board.Environment))
	if err != nil {
		log.Fatal(err)
	}
	if err := emulator.Run(profile); err != nil {
		log.Fatal(err)
	}
}
