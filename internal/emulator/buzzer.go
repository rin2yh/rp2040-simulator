//go:build !tinygo

package emulator

import (
	"bytes"
	"encoding/binary"
	"errors"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/rin2yh/rp2040-simulator/internal/bridge"
	"github.com/rin2yh/rp2040-simulator/internal/buzzer"
)

// Buzzer keeps RPC state separate from the audio player owned by the game loop.
type Buzzer struct {
	mu        sync.Mutex
	frequency uint32
	paused    bool
}

func (b *Buzzer) SetFrequency(args *bridge.BuzzerArgs, _ *struct{}) error {
	if b == nil {
		return errors.New("board has no buzzer")
	}
	if err := buzzer.ValidateFrequency(args.Frequency); err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.paused {
		return errors.New("buzzer is paused in bootloader")
	}
	b.frequency = args.Frequency
	return nil
}

func (b *Buzzer) Frequency() uint32 {
	if b == nil {
		return 0
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.frequency
}

func (b *Buzzer) reset(paused bool) {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.frequency, b.paused = 0, paused
}

// One second contains an integer number of cycles for integer Hz, so the
// stereo signed 16-bit PCM buffer loops without a phase discontinuity.
func buzzerPCM(hz uint32, sampleRate int) []byte {
	pcm := make([]byte, sampleRate*4)
	if hz == 0 {
		return pcm
	}
	for i := 0; i < sampleRate; i++ {
		sample := int16(2048)
		if (uint64(i)*uint64(hz)*2/uint64(sampleRate))%2 != 0 {
			sample = -sample
		}
		binary.LittleEndian.PutUint16(pcm[i*4:], uint16(sample))
		binary.LittleEndian.PutUint16(pcm[i*4+2:], uint16(sample))
	}
	return pcm
}

func (g *game) stopBuzzerAudio() {
	if g.buzzerPlayer != nil {
		_ = g.buzzerPlayer.Close()
		g.buzzerPlayer = nil
	}
	g.buzzerFrequency = 0
}

func (g *game) updateBuzzerAudio() error {
	hz := g.buzzer.Frequency()
	if hz == g.buzzerFrequency {
		return nil
	}
	g.stopBuzzerAudio()
	if hz == 0 {
		return nil
	}
	ctx := audio.CurrentContext()
	if ctx == nil {
		ctx = audio.NewContext(44100)
	}
	pcm := buzzerPCM(hz, ctx.SampleRate())
	player, err := ctx.NewPlayer(audio.NewInfiniteLoop(bytes.NewReader(pcm), int64(len(pcm))))
	if err != nil {
		return err
	}
	g.buzzerPlayer, g.buzzerFrequency = player, hz
	player.Play()
	return nil
}
