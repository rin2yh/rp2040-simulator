package emulator

import (
	"fmt"
	"sync"
)

// frameBuffer keeps drawing and presented values separate. Every operation
// touching either slice holds the same lock, including a complete frame write.
type frameBuffer[T any] struct {
	mu            sync.RWMutex
	buffer, front []T
}

func newFrameBuffer[T any](size int) frameBuffer[T] {
	return frameBuffer[T]{buffer: make([]T, size), front: make([]T, size)}
}

func (f *frameBuffer[T]) len() int { return len(f.buffer) }

func (f *frameBuffer[T]) update(index int, apply func(*T)) {
	f.mu.Lock()
	defer f.mu.Unlock()
	apply(&f.buffer[index])
}

func (f *frameBuffer[T]) presented(index int) T {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.front[index]
}

func (f *frameBuffer[T]) clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	clear(f.buffer)
}

func (f *frameBuffer[T]) display() {
	f.mu.Lock()
	defer f.mu.Unlock()
	copy(f.front, f.buffer)
}

func (f *frameBuffer[T]) write(values []T) error {
	if len(values) != f.len() {
		return fmt.Errorf("frame has %d values, want %d", len(values), f.len())
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	copy(f.buffer, values)
	copy(f.front, values)
	return nil
}
