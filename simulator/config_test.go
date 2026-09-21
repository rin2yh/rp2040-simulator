package simulator

import "testing"

func TestConfigure(t *testing.T) {
	type want struct {
		err bool
	}
	tests := []struct {
		name  string
		board Board
		want  want
	}{
		{name: "zero-kb02", board: BoardZeroKB02},
		{name: "conf2025badge", board: BoardConf2025Badge},
		{name: "empty", want: want{err: true}},
		{name: "unknown", board: Board("unknown"), want: want{err: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Configure(Config{Board: tt.board})
			if (err != nil) != tt.want.err {
				t.Fatalf("Configure(%q) error = %v, want error %v", tt.board, err, tt.want.err)
			}
		})
	}
}
