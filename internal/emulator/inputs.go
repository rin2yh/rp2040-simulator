package emulator

type Button struct{ Down bool }

func (b *Button) Pressed() bool { return b.Down }

type Joystick struct {
	X, Y float32
	Down bool
}

func (j *Joystick) Position() (float32, float32) { return j.X, j.Y }
func (j *Joystick) Pressed() bool                { return j.Down }
