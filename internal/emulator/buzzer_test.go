//go:build !tinygo

package emulator

import (
	"encoding/binary"
	"net"
	"net/rpc"
	"testing"

	"github.com/rin2yh/rp2040-simulator/internal/board"
	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

func TestBuzzerRPCAndReset(t *testing.T) {
	g, err := newGame(board.Conf2025Badge())
	if err != nil {
		t.Fatal(err)
	}
	server := rpc.NewServer()
	if err := server.RegisterName("Buzzer", g.buzzer); err != nil {
		t.Fatal(err)
	}
	a, b := net.Pipe()
	go server.ServeConn(a)
	client := rpc.NewClient(b)
	defer client.Close()
	set := func(hz uint32) error {
		return client.Call("Buzzer.SetFrequency", &bridge.BuzzerArgs{Frequency: hz}, &struct{}{})
	}
	if err := set(440); err != nil {
		t.Fatal(err)
	}
	if g.buzzer.Frequency() != 440 {
		t.Fatal("tone did not reach buzzer")
	}
	if err := set(20001); err == nil {
		t.Fatal("invalid frequency accepted")
	}
	if g.buzzer.Frequency() != 440 {
		t.Fatal("invalid frequency changed playing tone")
	}
	if err := set(0); err != nil {
		t.Fatal(err)
	}
	if g.buzzer.Frequency() != 0 {
		t.Fatal("Stop did not silence buzzer")
	}
	for _, bootloader := range []bool{false, true} {
		if err := set(880); err != nil {
			t.Fatal(err)
		}
		if err := g.restart(bootloader); err != nil {
			t.Fatal(err)
		}
		if g.buzzer.Frequency() != 0 {
			t.Fatal("tone survived reset")
		}
	}
	if err := set(440); err == nil {
		t.Fatal("bootloader accepted tone")
	}
	if err := g.restart(false); err != nil {
		t.Fatal(err)
	}
	if err := set(440); err != nil {
		t.Fatal(err)
	}
}

func TestBoardWithoutBuzzerRejectsRPC(t *testing.T) {
	server := rpc.NewServer()
	var buzzer *Buzzer
	if err := server.RegisterName("Buzzer", buzzer); err != nil {
		t.Fatal(err)
	}
	a, b := net.Pipe()
	go server.ServeConn(a)
	client := rpc.NewClient(b)
	defer client.Close()
	if err := client.Call("Buzzer.SetFrequency", &bridge.BuzzerArgs{Frequency: 440}, &struct{}{}); err == nil {
		t.Fatal("board without buzzer accepted tone")
	}
}

func TestBuzzerPCM(t *testing.T) {
	const rate = 44100
	for _, hz := range []uint32{0, 28, 440, 880, 20000} {
		pcm := buzzerPCM(hz, rate)
		if len(pcm) != rate*4 {
			t.Fatal("wrong PCM size")
		}
		crossings := 0
		previous := int16(binary.LittleEndian.Uint16(pcm[len(pcm)-4:]))
		for i := 0; i < len(pcm); i += 4 {
			left := int16(binary.LittleEndian.Uint16(pcm[i:]))
			right := int16(binary.LittleEndian.Uint16(pcm[i+2:]))
			if left != right {
				t.Fatal("stereo channels differ")
			}
			if hz == 0 && left != 0 {
				t.Fatal("silence is not zero")
			}
			if hz != 0 && left != 2048 && left != -2048 {
				t.Fatal("invalid square wave amplitude")
			}
			if left > 0 && previous < 0 {
				crossings++
			}
			previous = left
		}
		if crossings != int(hz) {
			t.Fatalf("%d Hz generated %d cycles", hz, crossings)
		}
	}
}
