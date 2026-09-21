package emulator

import (
	"image"
	"math"
)

type PointerLayout struct {
	Keys                          []image.Rectangle
	Encoder, Joystick             image.Point
	EncoderRadius, JoystickRadius int
}

type PointerInput struct {
	Position             image.Point
	Down, Right, Enabled bool
	Wheel                float64
}

// Pointer captures the control where a gesture starts. Dragging across other
// controls cannot activate them. Focus loss/modals cancel until mouse release.
type Pointer struct {
	Layout                    PointerLayout
	Keys                      []bool
	Delta                     float64
	EncoderDown, JoystickDown bool
	EncoderActive             bool
	X, Y                      float32
	capture                   int
	wasDown                   bool
	start, previous           image.Point
	dragged, angular          bool
}

const (
	captureNone     = -1
	captureEncoder  = -2
	captureJoystick = -3
)

func NewPointer(layout PointerLayout) *Pointer {
	return &Pointer{Layout: layout, Keys: make([]bool, len(layout.Keys)), capture: captureNone}
}

// Reset releases every input while retaining the board layout and buffers.
func (p *Pointer) Reset() {
	clear(p.Keys)
	*p = Pointer{Layout: p.Layout, Keys: p.Keys, capture: captureNone}
}

func InsideCircle(p, center image.Point, radius int) bool {
	d := p.Sub(center)
	return radius > 0 && d.X*d.X+d.Y*d.Y <= radius*radius
}

func (p *Pointer) Step(in PointerInput) {
	clear(p.Keys)
	p.Delta, p.X, p.Y = 0, 0, 0
	p.EncoderDown, p.JoystickDown = false, false
	p.EncoderActive = false
	if !in.Enabled {
		p.capture, p.wasDown = captureNone, in.Down
		return
	}
	if InsideCircle(in.Position, p.Layout.Encoder, p.Layout.EncoderRadius) {
		p.Delta = in.Wheel
		p.EncoderDown = in.Right
	}
	if InsideCircle(in.Position, p.Layout.Joystick, p.Layout.JoystickRadius) {
		p.JoystickDown = in.Right
	}
	if in.Down && !p.wasDown {
		p.start, p.previous, p.dragged = in.Position, in.Position, false
		p.capture = captureNone
		switch {
		case InsideCircle(in.Position, p.Layout.Encoder, p.Layout.EncoderRadius):
			p.capture = captureEncoder
			p.angular = !InsideCircle(in.Position, p.Layout.Encoder, p.Layout.EncoderRadius/2)
		case InsideCircle(in.Position, p.Layout.Joystick, p.Layout.JoystickRadius):
			p.capture = captureJoystick
		default:
			for i, rect := range p.Layout.Keys {
				if in.Position.In(rect) {
					p.capture = i
					break
				}
			}
		}
	}
	if in.Down && p.capture != captureNone {
		d := in.Position.Sub(p.start)
		if d.X*d.X+d.Y*d.Y >= 16 {
			p.dragged = true
		}
		switch p.capture {
		case captureEncoder:
			p.EncoderActive = true
			if !p.dragged {
				break
			}
			if !p.angular {
				p.Delta += float64(p.previous.Y-in.Position.Y) / 8
				break
			}
			a, b := p.previous.Sub(p.Layout.Encoder), in.Position.Sub(p.Layout.Encoder)
			// Ignore crossing the shaft, where the angle is undefined.
			if a.X*a.X+a.Y*a.Y <= 16 || b.X*b.X+b.Y*b.Y <= 16 {
				break
			}
			angle := math.Atan2(float64(b.Y), float64(b.X)) - math.Atan2(float64(a.Y), float64(a.X))
			p.Delta += math.Atan2(math.Sin(angle), math.Cos(angle)) * 20 / (2 * math.Pi)
		case captureJoystick:
			if p.dragged {
				x, y := float64(d.X)/float64(p.Layout.JoystickRadius), float64(d.Y)/float64(p.Layout.JoystickRadius)
				length := math.Max(1, math.Hypot(x, y))
				p.X, p.Y = float32(x/length), float32(y/length)
			}
		default:
			p.Keys[p.capture] = in.Position.In(p.Layout.Keys[p.capture])
		}
	}
	if !in.Down && p.wasDown {
		if !p.dragged {
			if p.capture == captureEncoder {
				p.EncoderDown = true
			}
			if p.capture == captureJoystick {
				p.JoystickDown = true
			}
		}
		p.capture = captureNone
	}
	p.previous, p.wasDown = in.Position, in.Down
}
