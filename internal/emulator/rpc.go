//go:build !tinygo

package emulator

import (
	"errors"
	"net"
	"net/rpc"

	"github.com/rin2yh/rp2040-simulator/internal/bridge"
)

type ledService struct {
	leds *LEDs
}

type inputService struct{ game *game }

func (s *inputService) Read(_ *struct{}, reply *bridge.Inputs) error {
	g := s.game
	g.mu.RLock()
	defer g.mu.RUnlock()
	if g.profile.Name != "zero-kb02" {
		return errors.New("board has no zero-kb02 controls")
	}
	if g.bootloader {
		return errors.New("application paused in bootloader")
	}
	for i := range reply.Keys {
		reply.Keys[i] = g.buttons[i].Down
	}
	reply.EncoderPosition = g.encoder.Position()
	reply.EncoderPressed = g.encoder.Pressed()
	reply.JoystickX, reply.JoystickY = g.joystick.Position()
	reply.JoystickPressed = g.joystick.Pressed()
	reply.Generation = g.generation
	return nil
}

type displayService struct{ display *Display }

func (s *displayService) WriteFrame(args *bridge.DisplayFrame, _ *struct{}) error {
	return s.display.WriteFrame(args.Pixels)
}

func (s *ledService) WriteColors(args *bridge.WriteLEDsArgs, _ *struct{}) error {
	if s.leds == nil {
		return errors.New("board has no LEDs")
	}
	return s.leds.WriteColors(args.Colors)
}

func startRPC(g *game) (net.Listener, error) {
	server := rpc.NewServer()
	if err := server.RegisterName("LEDs", &ledService{leds: g.leds}); err != nil {
		return nil, err
	}
	if err := server.RegisterName("Buzzer", g.buzzer); err != nil {
		return nil, err
	}
	if err := server.RegisterName("Inputs", &inputService{game: g}); err != nil {
		return nil, err
	}
	if err := server.RegisterName("Display", &displayService{display: g.display}); err != nil {
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
