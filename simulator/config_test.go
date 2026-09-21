package simulator

import "testing"

func TestConfigureRejectsInvalidBoard(t *testing.T) {
	for _, board := range []Board{"", Board("unknown")} {
		if err := Configure(Config{Board: board}); err == nil {
			t.Fatalf("Configure(%q) accepted an invalid board", board)
		}
	}
}

func TestConfigureAcceptsSupportedBoards(t *testing.T) {
	for _, board := range []Board{BoardZeroKB02, BoardConf2025Badge} {
		if err := Configure(Config{Board: board}); err != nil {
			t.Fatalf("Configure(%q): %v", board, err)
		}
	}
}
