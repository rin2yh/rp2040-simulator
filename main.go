package main

import (
	"flag"
	"log"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/emulator"
)

func main() {
	name := flag.String("board", "zero-kb02", "board to emulate")
	flag.Parse()
	profile, err := board.Lookup(*name)
	if err != nil {
		log.Fatal(err)
	}
	if err := emulator.Run(profile); err != nil {
		log.Fatal(err)
	}
}
