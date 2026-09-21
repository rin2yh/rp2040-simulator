package emulator

// Encoder stores detents and accumulates fractional trackpad wheel events.
type Encoder struct {
	position  int
	pressed   bool
	remainder float64
}

func (e *Encoder) Position() int           { return e.position }
func (e *Encoder) Pressed() bool           { return e.pressed }
func (e *Encoder) SetPressed(pressed bool) { e.pressed = pressed }

func (e *Encoder) Rotate(detents float64) {
	e.remainder += detents
	whole := int(e.remainder)
	e.position += whole
	e.remainder -= float64(whole)
}
