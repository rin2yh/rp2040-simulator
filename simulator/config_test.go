package simulator

import (
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/tabletest"
)

func TestConfigure(t *testing.T) {
	type want struct {
		err bool
	}
	cases := []tabletest.Case[Board, want]{
		{Name: "zero-kb02", In: BoardZeroKB02},
		{Name: "conf2025badge", In: BoardConf2025Badge},
		{Name: "empty", Want: want{err: true}},
		{Name: "unknown", In: Board("unknown"), Want: want{err: true}},
	}

	tabletest.Run(t, cases, func(t *testing.T, board Board, want want) {
		err := Configure(Config{Board: board})
		if (err != nil) != want.err {
			t.Fatalf("Configure(%q) error = %v, want error %v", board, err, want.err)
		}
	})
}
