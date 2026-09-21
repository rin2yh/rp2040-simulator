package simulator

import "testing"

func TestConfigureRejectsUnknownBoard(t *testing.T) {
	err := Configure(Config{Board: "unknown"})
	if err == nil {
		t.Fatal("Configure accepted an unknown board")
	}
}

func TestConfigureAcceptsSupportedBoards(t *testing.T) {
	for _, name := range []string{"", "zero-kb02", "conf2025badge"} {
		if err := Configure(Config{Board: name}); err != nil {
			t.Fatalf("Configure(%q): %v", name, err)
		}
	}
}
