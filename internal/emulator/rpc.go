//go:build !tinygo

package emulator

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

type ledService struct {
	leds *LEDs
}

func (s *ledService) WriteColors(args *bridge.WriteLEDsArgs, _ *struct{}) error {
	if s.leds == nil {
		return errors.New("board has no LEDs")
	}
	if len(args.Colors) != s.leds.Len() {
		return fmt.Errorf("got %d LED colors, want %d", len(args.Colors), s.leds.Len())
	}
	for i, c := range args.Colors {
		s.leds.Set(i, c)
	}
	return s.leds.Display()
}

func startRPC(leds *LEDs) (net.Listener, error) {
	server := rpc.NewServer()
	if err := server.RegisterName("LEDs", &ledService{leds: leds}); err != nil {
		return nil, err
	}
	listener, err := net.Listen("tcp", bridge.Address)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()
	return listener, nil
}
